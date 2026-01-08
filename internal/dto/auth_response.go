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
