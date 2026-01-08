package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	jwtPkg "github.com/AnggaKay/ojek-kampus-backend/pkg/jwt"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/password"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/utils"
)

type AuthService interface {
	RegisterPassenger(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *authService) RegisterPassenger(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	// Check if phone already exists
	existing, _ := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if existing != nil {
		return nil, fmt.Errorf("phone number already registered")
	}

	// Check email if provided
	if req.Email != nil && *req.Email != "" {
		existing, _ := s.userRepo.FindByEmail(ctx, *req.Email)
		if existing != nil {
			return nil, fmt.Errorf("email already registered")
		}
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		PhoneNumber:   phoneNumber,
		PasswordHash:  hashedPassword,
		Email:         req.Email,
		FullName:      req.FullName,
		Role:          entity.RolePassenger,
		Status:        entity.StatusActive,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID, string(user.Role), req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User: &dto.UserResponse{
			ID:            user.ID,
			PhoneNumber:   user.PhoneNumber,
			Email:         user.Email,
			FullName:      user.FullName,
			Role:          user.Role,
			Status:        user.Status,
			PhoneVerified: user.PhoneVerified,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	// Find user
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number or password")
	}

	// Verify password
	if !password.Verify(user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid phone number or password")
	}

	// Check if account is suspended
	if user.Status == entity.StatusSuspended {
		return nil, fmt.Errorf("account is suspended")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
	}

	// Generate tokens
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID, string(user.Role), req.DeviceInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User: &dto.UserResponse{
			ID:            user.ID,
			PhoneNumber:   user.PhoneNumber,
			Email:         user.Email,
			FullName:      user.FullName,
			Role:          user.Role,
			Status:        user.Status,
			PhoneVerified: user.PhoneVerified,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	// Hash the refresh token
	tokenHash := hashToken(refreshToken)

	// Find token in database
	token, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if revoked
	if token.IsRevoked {
		return nil, fmt.Errorf("token has been revoked")
	}

	// Check if expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Update last used
	if err := s.refreshTokenRepo.UpdateLastUsed(ctx, token.ID); err != nil {
		// Log error but don't fail
	}

	// Find user
	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, token.UserType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)

	token, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	return s.refreshTokenRepo.Revoke(ctx, token.ID, "LOGOUT")
}

// Helper functions
func (s *authService) createRefreshToken(ctx context.Context, userID int, userType, deviceInfo string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	tokenHash := hashToken(token)

	// Create token record
	refreshToken := &entity.RefreshToken{
		UserID:     userID,
		UserType:   userType,
		TokenHash:  tokenHash,
		DeviceInfo: &deviceInfo,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		IsRevoked:  false,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	return token, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
