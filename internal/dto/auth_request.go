package dto

// Auth Requests
type RegisterRequest struct {
	PhoneNumber string  `json:"phone_number" validate:"required,min=10,max=15"`
	Password    string  `json:"password" validate:"required,min=8"`
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
	DeviceInfo  string `json:"device_info,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
