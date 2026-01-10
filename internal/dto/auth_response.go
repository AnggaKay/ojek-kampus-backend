package dto

import "github.com/AnggaKay/ojek-kampus-backend/internal/entity"

// Auth Responses
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"` // seconds
}

type UserResponse struct {
	ID            int               `json:"id"`
	PhoneNumber   string            `json:"phone_number"`
	Email         *string           `json:"email,omitempty"`
	FullName      string            `json:"full_name"`
	Role          entity.UserRole   `json:"role"`
	Status        entity.UserStatus `json:"status"`
	PhoneVerified bool              `json:"phone_verified"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// LoginResponse is a unified response for login that includes role-specific profile
type LoginResponse struct {
	User             *UserResponse             `json:"user"`
	PassengerProfile *PassengerProfileResponse `json:"passenger_profile,omitempty"` // Only for PASSENGER role
	DriverProfile    *DriverProfileResponse    `json:"driver_profile,omitempty"`    // Only for DRIVER role
	AccessToken      string                    `json:"access_token"`
	RefreshToken     string                    `json:"refresh_token"`
	ExpiresIn        int                       `json:"expires_in"` // seconds
}

// PassengerProfileResponse represents passenger profile data
type PassengerProfileResponse struct {
	ID                    int     `json:"id"`
	UserID                int     `json:"user_id"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	HomeAddress           *string `json:"home_address,omitempty"`
	TotalCompletedOrders  int     `json:"total_completed_orders"`
}
