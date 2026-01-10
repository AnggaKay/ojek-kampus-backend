package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	jwtPkg "github.com/AnggaKay/ojek-kampus-backend/pkg/jwt"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/password"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/storage"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/utils"
)

type DriverService interface {
	RegisterDriver(ctx context.Context, req dto.RegisterDriverRequest, files map[string]*multipart.FileHeader) (*dto.DriverAuthResponse, error)
}

type driverService struct {
	userRepo         repository.UserRepository
	driverRepo       repository.DriverRepository
	refreshTokenRepo repository.RefreshTokenRepository
	fileStorage      storage.FileStorage
	tokenHelper      *TokenHelper
}

func NewDriverService(
	userRepo repository.UserRepository,
	driverRepo repository.DriverRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	fileStorage storage.FileStorage,
) DriverService {
	return &driverService{
		userRepo:         userRepo,
		driverRepo:       driverRepo,
		refreshTokenRepo: refreshTokenRepo,
		fileStorage:      fileStorage,
		tokenHelper:      NewTokenHelper(refreshTokenRepo),
	}
}

func (s *driverService) RegisterDriver(
	ctx context.Context,
	req dto.RegisterDriverRequest,
	files map[string]*multipart.FileHeader,
) (*dto.DriverAuthResponse, error) {
	logger.Log.Info().Str("phone", req.PhoneNumber).Msg("Driver registration attempt")

	// 1. Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		logger.Log.Warn().Err(err).Str("phone", req.PhoneNumber).Msg("Password validation failed")
		return nil, err
	}

	// 2. Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	// 3. Check if phone already exists
	phoneExists, err := s.userRepo.ExistsByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		logger.Log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to check phone existence")
		return nil, fmt.Errorf("failed to check phone availability")
	}
	if phoneExists {
		logger.Log.Warn().Str("phone", phoneNumber).Msg("Phone number already registered")
		return nil, fmt.Errorf(constants.ErrPhoneAlreadyRegistered)
	}

	// 4. Check if email already exists (if provided)
	if req.Email != "" {
		emailExists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			logger.Log.Error().Err(err).Str("email", req.Email).Msg("Failed to check email existence")
			return nil, fmt.Errorf("failed to check email availability")
		}
		if emailExists {
			logger.Log.Warn().Str("email", req.Email).Msg("Email already registered")
			return nil, fmt.Errorf(constants.ErrEmailAlreadyRegistered)
		}
	}

	// 5. Check if vehicle plate already exists
	plateExists, err := s.driverRepo.ExistsByVehiclePlate(ctx, req.VehiclePlate)
	if err != nil {
		logger.Log.Error().Err(err).Str("plate", req.VehiclePlate).Msg("Failed to check vehicle plate existence")
		return nil, fmt.Errorf("failed to check vehicle plate availability")
	}
	if plateExists {
		logger.Log.Warn().Str("plate", req.VehiclePlate).Msg("Vehicle plate already registered")
		return nil, fmt.Errorf(constants.ErrVehiclePlateExists)
	}

	// 6. Validate all required files are present
	requiredFiles := []string{"ktp", "sim", "stnk", "ktm"}
	for _, docType := range requiredFiles {
		if files[docType] == nil {
			logger.Log.Warn().Str("doc_type", docType).Msg("Required document missing")
			return nil, fmt.Errorf("document %s is required", docType)
		}
	}

	// 7. Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf(constants.ErrFailedToHashPassword+": %w", err)
	}

	// 8. Create user (role=DRIVER)
	var emailPtr *string
	if req.Email != "" {
		emailPtr = &req.Email
	}

	user := &entity.User{
		PhoneNumber:   phoneNumber,
		PasswordHash:  hashedPassword,
		Email:         emailPtr,
		FullName:      req.FullName,
		Role:          entity.RoleDriver,
		Status:        entity.StatusActive,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Log.Error().Err(err).Str("phone", phoneNumber).Msg("Failed to create user")
		return nil, fmt.Errorf(constants.ErrFailedToCreateUser+": %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Str("phone", phoneNumber).Msg("Driver user created successfully")

	// 9. Upload documents
	uploadedDocs := make(map[string]string)
	var uploadErr error

	for docType, file := range files {
		if file == nil {
			continue
		}

		filePath, err := s.fileStorage.Upload(file, user.ID, docType)
		if err != nil {
			logger.Log.Error().Err(err).Int("user_id", user.ID).Str("doc_type", docType).Msg("Failed to upload document")
			uploadErr = err
			break
		}
		uploadedDocs[docType] = filePath
		logger.Log.Info().Int("user_id", user.ID).Str("doc_type", docType).Str("path", filePath).Msg("Document uploaded")
	}

	// Rollback: delete uploaded files if any error occurred
	if uploadErr != nil {
		logger.Log.Warn().Int("user_id", user.ID).Msg("Rolling back uploaded files")
		for _, filePath := range uploadedDocs {
			_ = s.fileStorage.Delete(filePath)
		}
		return nil, fmt.Errorf(constants.ErrFailedToUploadFile+": %w", uploadErr)
	}

	// 10. Create driver profile
	driverProfile := &entity.DriverProfile{
		UserID:       user.ID,
		VehicleType:  constants.VehicleTypeMotor, // Always MOTOR as specified
		VehiclePlate: req.VehiclePlate,
		IsVerified:   false,
		IsActive:     false,
	}

	// Set optional fields
	if req.VehicleBrand != "" {
		driverProfile.VehicleBrand = &req.VehicleBrand
	}
	if req.VehicleModel != "" {
		driverProfile.VehicleModel = &req.VehicleModel
	}
	if req.VehicleColor != "" {
		driverProfile.VehicleColor = &req.VehicleColor
	}

	// Set document paths
	if ktpPath, ok := uploadedDocs["ktp"]; ok {
		driverProfile.KTPPhoto = &ktpPath
	}
	if simPath, ok := uploadedDocs["sim"]; ok {
		driverProfile.SIMPhoto = &simPath
	}
	if stnkPath, ok := uploadedDocs["stnk"]; ok {
		driverProfile.STNKPhoto = &stnkPath
	}
	if ktmPath, ok := uploadedDocs["ktm"]; ok {
		driverProfile.KTMPhoto = &ktmPath
	}

	if err := s.driverRepo.Create(ctx, driverProfile); err != nil {
		logger.Log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to create driver profile")
		// Rollback: delete uploaded files
		for _, filePath := range uploadedDocs {
			_ = s.fileStorage.Delete(filePath)
		}
		return nil, fmt.Errorf("failed to create driver profile: %w", err)
	}

	logger.Log.Info().Int("user_id", user.ID).Int("profile_id", driverProfile.ID).Msg("Driver profile created successfully")

	// 11. Generate tokens
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

	logger.Log.Info().Int("user_id", user.ID).Msg("Driver registration completed successfully")

	// 12. Build response
	return &dto.DriverAuthResponse{
		User: &dto.UserResponse{
			ID:            user.ID,
			PhoneNumber:   user.PhoneNumber,
			Email:         user.Email,
			FullName:      user.FullName,
			Role:          user.Role,
			Status:        user.Status,
			PhoneVerified: user.PhoneVerified,
		},
		DriverProfile: &dto.DriverProfileResponse{
			ID:                   driverProfile.ID,
			UserID:               driverProfile.UserID,
			VehicleType:          driverProfile.VehicleType,
			VehiclePlate:         driverProfile.VehiclePlate,
			VehicleBrand:         driverProfile.VehicleBrand,
			VehicleModel:         driverProfile.VehicleModel,
			VehicleColor:         driverProfile.VehicleColor,
			IsVerified:           driverProfile.IsVerified,
			VerificationStatus:   constants.VerificationStatusPending,
			IsActive:             driverProfile.IsActive,
			TotalCompletedOrders: driverProfile.TotalCompletedOrders,
			RatingAvg:            driverProfile.RatingAvg,
			Documents: &dto.Documents{
				KTPUploaded:  driverProfile.KTPPhoto != nil,
				SIMUploaded:  driverProfile.SIMPhoto != nil,
				STNKUploaded: driverProfile.STNKPhoto != nil,
				KTMUploaded:  driverProfile.KTMPhoto != nil,
			},
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(constants.AccessTokenTTL.Seconds()),
	}, nil
}

// createRefreshToken creates a refresh token
func (s *driverService) createRefreshToken(ctx context.Context, userID int, userType, deviceInfo string) (string, error) {
	return s.tokenHelper.CreateRefreshToken(ctx, userID, userType, deviceInfo)
}
