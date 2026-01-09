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
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	jwtPkg "github.com/AnggaKay/ojek-kampus-backend/pkg/jwt"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
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
	logger.Log.Info().Str("phone", req.PhoneNumber).Msg("Passenger registration attempt")

	// Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		logger.Log.Warn().Err(err).Str("phone", req.PhoneNumber).Msg("Password validation failed")
		return nil, err
	}

	// Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	// Check if phone already exists (using optimized Exists method)
	phoneExists, err := s.userRepo.ExistsByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		logger.Log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to check phone existence")
		return nil, fmt.Errorf("failed to check phone availability")
	}
	if phoneExists {
		logger.Log.Warn().Str("phone", phoneNumber).Msg("Phone number already registered")
		return nil, fmt.Errorf(constants.ErrPhoneAlreadyRegistered)
	}

	// Check email if provided
	if req.Email != nil && *req.Email != "" {
		emailExists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			logger.Log.Error().Err(err).Str("email", *req.Email).Msg("Failed to check email existence")
			return nil, fmt.Errorf("failed to check email availability")
		}
		if emailExists {
			logger.Log.Warn().Str("email", *req.Email).Msg("Email already registered")
			return nil, fmt.Errorf(constants.ErrEmailAlreadyRegistered)
		}
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf(constants.ErrFailedToHashPassword+": %w", err)
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
		logger.Log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to create user")
		return nil, fmt.Errorf(constants.ErrFailedToCreateUser+": %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Str("phone", phoneNumber).Msg("User created successfully")

	// Generate tokens
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, string(user.Role))
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate access token")
		return nil, fmt.Errorf(constants.ErrFailedToGenerateToken+": %w", err)
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID, string(user.Role), req.PhoneNumber)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate refresh token")
		return nil, fmt.Errorf(constants.ErrFailedToGenerateToken+": %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Msg("Registration completed successfully")

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
		ExpiresIn:    int(constants.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	logger.Log.Info().Str("phone", phoneNumber).Msg("Login attempt")

	// Find user
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		logger.Log.Warn().Str("phone", phoneNumber).Msg("Login failed: user not found")
		return nil, fmt.Errorf(constants.ErrInvalidCredentials)
	}

	// Verify password
	if !password.Verify(user.PasswordHash, req.Password) {
		logger.Log.Warn().Int("user_id", user.ID).Str("phone", phoneNumber).Msg("Login failed: invalid password")
		return nil, fmt.Errorf(constants.ErrInvalidCredentials)
	}

	// Check if account is suspended
	if user.Status == entity.StatusSuspended {
		logger.Log.Warn().Int("user_id", user.ID).Msg("Login attempt on suspended account")
		return nil, fmt.Errorf(constants.ErrAccountSuspended)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to update last login")
		// Don't fail login for this
	}

	// Generate tokens
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, string(user.Role))
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate access token")
		return nil, fmt.Errorf(constants.ErrFailedToGenerateToken+": %w", err)
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID, string(user.Role), req.DeviceInfo)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate refresh token")
		return nil, fmt.Errorf(constants.ErrFailedToGenerateToken+": %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Str("phone", phoneNumber).Msg("Login successful")

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
		ExpiresIn:    int(constants.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	logger.Log.Debug().Msg("Refresh token attempt")

	// Hash the refresh token
	tokenHash := hashToken(refreshToken)

	// Find token in database
	token, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		logger.Log.Warn().Msg("Invalid refresh token provided")
		return nil, fmt.Errorf(constants.ErrInvalidRefreshToken)
	}

	// Check if revoked
	if token.IsRevoked {
		logger.Log.Warn().Int("token_id", token.ID).Int("user_id", token.UserID).Msg("Attempted to use revoked token")
		return nil, fmt.Errorf(constants.ErrTokenRevoked)
	}

	// Check if expired
	if token.ExpiresAt.Before(time.Now()) {
		logger.Log.Warn().Int("token_id", token.ID).Int("user_id", token.UserID).Msg("Attempted to use expired token")
		return nil, fmt.Errorf(constants.ErrTokenExpired)
	}

	// Update last used
	if err := s.refreshTokenRepo.UpdateLastUsed(ctx, token.ID); err != nil {
		logger.Log.Error().Err(err).Int("token_id", token.ID).Msg("Failed to update token last used")
		// Don't fail refresh for this
	}

	// Find user
	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", token.UserID).Msg("User not found for valid token")
		return nil, fmt.Errorf(constants.ErrUserNotFound)
	}

	// Generate new access token
	accessToken, err := jwtPkg.GenerateAccessToken(user.ID, user.Role, token.UserType)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate access token")
		return nil, fmt.Errorf(constants.ErrFailedToGenerateToken+": %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Msg("Token refreshed successfully")

	return &dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(constants.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	logger.Log.Debug().Msg("Logout attempt")

	tokenHash := hashToken(refreshToken)

	token, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		logger.Log.Warn().Msg("Invalid refresh token on logout")
		return fmt.Errorf(constants.ErrInvalidRefreshToken)
	}

	if err := s.refreshTokenRepo.Revoke(ctx, token.ID, constants.RevokeReasonLogout); err != nil {
		logger.Log.Error().Err(err).Int("token_id", token.ID).Msg("Failed to revoke token")
		return err
	}

	logger.Log.Info().Int("user_id", token.UserID).Int("token_id", token.ID).Msg("Logout successful")
	return nil
}

// Helper functions
func (s *authService) createRefreshToken(ctx context.Context, userID int, userType, deviceInfo string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, constants.RefreshTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate random bytes for refresh token")
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
		ExpiresAt:  time.Now().Add(constants.RefreshTokenTTL),
		IsRevoked:  false,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("Failed to save refresh token")
		return "", err
	}

	return token, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
