package dto

// ============================================================================
// Common Request DTOs
// ============================================================================

// LoginRequest represents login credentials for all user types
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
	DeviceInfo  string `json:"device_info,omitempty"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ============================================================================
// Common Response DTOs
// ============================================================================

// UserResponse represents basic user information
type UserResponse struct {
	ID            int     `json:"id"`
	PhoneNumber   string  `json:"phone_number"`
	Email         *string `json:"email,omitempty"`
	FullName      string  `json:"full_name"`
	Role          string  `json:"role"`   // Will be converted from entity.UserRole to string
	Status        string  `json:"status"` // Will be converted from entity.UserStatus to string
	PhoneVerified bool    `json:"phone_verified"`
}

// TokenResponse represents token refresh response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

// Documents represents document upload status
type Documents struct {
	KTPUploaded  bool `json:"ktp_uploaded"`
	SIMUploaded  bool `json:"sim_uploaded"`
	STNKUploaded bool `json:"stnk_uploaded"`
	KTMUploaded  bool `json:"ktm_uploaded"`
}

// ============================================================================
// Authentication Response DTOs
// ============================================================================

// AuthResponse represents basic authentication response (used for registration)
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"` // seconds
}

// LoginResponse represents unified login response with role-specific profile
type LoginResponse struct {
	User             *UserResponse             `json:"user"`
	PassengerProfile *PassengerProfileResponse `json:"passenger_profile,omitempty"` // Only for PASSENGER
	DriverProfile    *DriverProfileResponse    `json:"driver_profile,omitempty"`    // Only for DRIVER
	AccessToken      string                    `json:"access_token"`
	RefreshToken     string                    `json:"refresh_token"`
	ExpiresIn        int                       `json:"expires_in"` // seconds
}
