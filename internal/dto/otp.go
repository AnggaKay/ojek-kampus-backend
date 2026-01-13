package dto

import "github.com/AnggaKay/ojek-kampus-backend/internal/entity"

// SendOTPRequest represents request to send OTP code
type SendOTPRequest struct {
	PhoneNumber string            `json:"phone_number" validate:"required,e164"`
	Purpose     entity.OTPPurpose `json:"purpose" validate:"required,oneof=REGISTRATION PASSWORD_RESET PHONE_VERIFICATION"`
}

// VerifyOTPRequest represents request to verify OTP code
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	OTPCode     string `json:"otp_code" validate:"required,len=6,numeric"`
}

// ResendOTPRequest represents request to resend OTP code
type ResendOTPRequest struct {
	PhoneNumber string            `json:"phone_number" validate:"required,e164"`
	Purpose     entity.OTPPurpose `json:"purpose" validate:"required,oneof=REGISTRATION PASSWORD_RESET PHONE_VERIFICATION"`
}

// SendOTPResponse represents response after sending OTP
type SendOTPResponse struct {
	PhoneNumber string `json:"phone_number"`
	ExpiresIn   int    `json:"expires_in"` // seconds
	Message     string `json:"message"`
}

// VerifyOTPResponse represents response after successful OTP verification
type VerifyOTPResponse struct {
	PhoneNumber       string `json:"phone_number"`
	Verified          bool   `json:"verified"`
	VerificationToken string `json:"verification_token,omitempty"` // JWT token for next step
	Message           string `json:"message"`
}

// OTPErrorResponse represents detailed error response for OTP operations
type OTPErrorResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	RemainingTime *int   `json:"remaining_time,omitempty"` // seconds until can resend
	AttemptsLeft  *int   `json:"attempts_left,omitempty"`  // verification attempts left
}
