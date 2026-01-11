package mapper

import (
	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
)

// ============================================================================
// User Mappers
// ============================================================================

// ToUserResponse converts entity.User to dto.UserResponse
func ToUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            user.ID,
		PhoneNumber:   user.PhoneNumber,
		Email:         user.Email,
		FullName:      user.FullName,
		Role:          string(user.Role),
		Status:        string(user.Status),
		PhoneVerified: user.PhoneVerified,
	}
}

// ============================================================================
// Passenger Profile Mappers
// ============================================================================

// ToPassengerProfileResponse converts entity.PassengerProfile to dto.PassengerProfileResponse
func ToPassengerProfileResponse(profile *entity.PassengerProfile) *dto.PassengerProfileResponse {
	if profile == nil {
		return nil
	}

	return &dto.PassengerProfileResponse{
		ID:                    profile.ID,
		UserID:                profile.UserID,
		EmergencyContactName:  profile.EmergencyContactName,
		EmergencyContactPhone: profile.EmergencyContactPhone,
		HomeAddress:           profile.HomeAddress,
		TotalCompletedOrders:  profile.TotalCompletedOrders,
	}
}

// ============================================================================
// Driver Profile Mappers
// ============================================================================

// ToDriverProfileResponse converts entity.DriverProfile to dto.DriverProfileResponse
func ToDriverProfileResponse(profile *entity.DriverProfile) *dto.DriverProfileResponse {
	if profile == nil {
		return nil
	}

	// Determine verification status
	verificationStatus := determineVerificationStatus(profile)

	return &dto.DriverProfileResponse{
		ID:                   profile.ID,
		UserID:               profile.UserID,
		VehicleType:          profile.VehicleType,
		VehiclePlate:         profile.VehiclePlate,
		VehicleBrand:         profile.VehicleBrand,
		VehicleModel:         profile.VehicleModel,
		VehicleColor:         profile.VehicleColor,
		IsVerified:           profile.IsVerified,
		VerificationStatus:   verificationStatus,
		RejectionReason:      profile.RejectionReason,
		IsActive:             profile.IsActive,
		TotalCompletedOrders: profile.TotalCompletedOrders,
		RatingAvg:            profile.RatingAvg,
		Documents:            toDocumentsResponse(profile),
	}
}

// determineVerificationStatus determines the verification status based on profile data
func determineVerificationStatus(profile *entity.DriverProfile) string {
	if profile.IsVerified {
		return constants.VerificationStatusVerified
	}
	if profile.RejectionReason != nil && *profile.RejectionReason != "" {
		return constants.VerificationStatusRejected
	}
	return constants.VerificationStatusPending
}

// toDocumentsResponse converts driver document paths to Documents DTO
func toDocumentsResponse(profile *entity.DriverProfile) *dto.Documents {
	return &dto.Documents{
		KTPUploaded:  profile.KTPPhoto != nil && *profile.KTPPhoto != "",
		SIMUploaded:  profile.SIMPhoto != nil && *profile.SIMPhoto != "",
		STNKUploaded: profile.STNKPhoto != nil && *profile.STNKPhoto != "",
		KTMUploaded:  profile.KTMPhoto != nil && *profile.KTMPhoto != "",
	}
}

// ============================================================================
// Auth Response Builders
// ============================================================================

// BuildAuthResponse builds standard auth response (for registration)
func BuildAuthResponse(user *entity.User, accessToken, refreshToken string, expiresIn int) *dto.AuthResponse {
	return &dto.AuthResponse{
		User:         ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

// BuildDriverAuthResponse builds driver auth response (for driver registration)
func BuildDriverAuthResponse(
	user *entity.User,
	driverProfile *entity.DriverProfile,
	accessToken,
	refreshToken string,
	expiresIn int,
) *dto.DriverAuthResponse {
	return &dto.DriverAuthResponse{
		User:          ToUserResponse(user),
		DriverProfile: ToDriverProfileResponse(driverProfile),
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ExpiresIn:     expiresIn,
	}
}

// BuildLoginResponse builds unified login response with role-specific profile
func BuildLoginResponse(
	user *entity.User,
	accessToken,
	refreshToken string,
	expiresIn int,
	passengerProfile *entity.PassengerProfile,
	driverProfile *entity.DriverProfile,
) *dto.LoginResponse {
	response := &dto.LoginResponse{
		User:         ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	// Add role-specific profile
	if passengerProfile != nil {
		response.PassengerProfile = ToPassengerProfileResponse(passengerProfile)
	}
	if driverProfile != nil {
		response.DriverProfile = ToDriverProfileResponse(driverProfile)
	}

	return response
}
