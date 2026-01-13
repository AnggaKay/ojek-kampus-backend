package entity

import "time"

// OTPPurpose defines the purpose of OTP code
type OTPPurpose string

const (
	OTPPurposeRegistration      OTPPurpose = "REGISTRATION"
	OTPPurposePasswordReset     OTPPurpose = "PASSWORD_RESET"
	OTPPurposePhoneVerification OTPPurpose = "PHONE_VERIFICATION"
)

// OTPCode represents the otp_codes table
type OTPCode struct {
	ID          int        `json:"id" db:"id"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	OTPCode     string     `json:"otp_code" db:"otp_code"`
	Purpose     OTPPurpose `json:"purpose" db:"purpose"`
	ExpiresAt   time.Time  `json:"expires_at" db:"expires_at"`
	IsUsed      bool       `json:"is_used" db:"is_used"`
	UsedAt      *time.Time `json:"used_at,omitempty" db:"used_at"`
	Attempts    int        `json:"attempts" db:"attempts"`
	IPAddress   *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent   *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// IsExpired checks if the OTP code has expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid checks if OTP is valid (not used and not expired)
func (o *OTPCode) IsValid() bool {
	return !o.IsUsed && !o.IsExpired()
}

// CanRetry checks if user can retry verification (max 3 attempts)
func (o *OTPCode) CanRetry() bool {
	return o.Attempts < 3
}
