package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/internal/mapper"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	jwtPkg "github.com/AnggaKay/ojek-kampus-backend/pkg/jwt"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/password"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/utils"
)

type AuthService interface {
	RegisterPassenger(ctx context.Context, req dto.RegisterPassengerRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	userRepo         repository.UserRepository
	passengerRepo    repository.PassengerRepository
	driverRepo       repository.DriverRepository
	refreshTokenRepo repository.RefreshTokenRepository
	tokenHelper      *TokenHelper
}

func NewAuthService(
	userRepo repository.UserRepository,
	passengerRepo repository.PassengerRepository,
	driverRepo repository.DriverRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		passengerRepo:    passengerRepo,
		driverRepo:       driverRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenHelper:      NewTokenHelper(refreshTokenRepo),
	}
}

func (s *authService) RegisterPassenger(ctx context.Context, req dto.RegisterPassengerRequest) (*dto.AuthResponse, error) {
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

	// Create passenger profile
	passengerProfile := &entity.PassengerProfile{
		UserID:             user.ID,
		ProfilePicture:     nil,
		FCMToken:           nil,
		TotalOrders:        0,
		TotalCancellations: 0,
		RatingAvg:          0.0,
	}

	if err := s.passengerRepo.Create(ctx, passengerProfile); err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to create passenger profile")
		// Note: User already created, but profile creation failed
		// In production, consider transaction rollback
		return nil, fmt.Errorf("failed to create passenger profile: %w", err)
	}

	logger.Log.Info().Int("passenger_id", passengerProfile.ID).Int("user_id", user.ID).Msg("Passenger profile created successfully")

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

	return mapper.BuildAuthResponse(user, accessToken, refreshToken, int(constants.AccessTokenTTL.Seconds())), nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
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

	// Prepare base response
	var driverProfile *entity.DriverProfile
	var passengerProfile *entity.PassengerProfile

	// Fetch role-specific profile
	if user.Role == entity.RoleDriver {
		driverProfile, _ = s.fetchDriverProfile(ctx, user.ID)
	} else if user.Role == entity.RolePassenger {
		passengerProfile, _ = s.fetchPassengerProfile(ctx, user.ID)
	}

	logger.Log.Info().Int("user_id", user.ID).Str("phone", phoneNumber).Str("role", string(user.Role)).Msg("Login successful")

	return mapper.BuildLoginResponse(
		user,
		accessToken,
		refreshToken,
		int(constants.AccessTokenTTL.Seconds()),
		passengerProfile,
		driverProfile,
	), nil
}

// fetchDriverProfile fetches driver profile by user ID
func (s *authService) fetchDriverProfile(ctx context.Context, userID int) (*entity.DriverProfile, error) {
	driverProfile, err := s.driverRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("Failed to fetch driver profile")
		return nil, err
	}
	logger.Log.Info().Int("driver_id", driverProfile.ID).Msg("Driver profile fetched")
	return driverProfile, nil
}

// fetchPassengerProfile fetches passenger profile by user ID
func (s *authService) fetchPassengerProfile(ctx context.Context, userID int) (*entity.PassengerProfile, error) {
	passengerProfile, err := s.passengerRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("Failed to fetch passenger profile")
		return nil, err
	}
	logger.Log.Info().Int("passenger_id", passengerProfile.ID).Msg("Passenger profile fetched")
	return passengerProfile, nil
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

// createRefreshToken creates a refresh token
func (s *authService) createRefreshToken(ctx context.Context, userID int, userType, deviceInfo string) (string, error) {
	return s.tokenHelper.CreateRefreshToken(ctx, userID, userType, deviceInfo)
}

func hashToken(token string) string {
	return HashToken(token)
}
