package constants

import "time"

// Authentication constants
const (
	// JWT Token Time-to-Live
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour

	// Password requirements
	MinPasswordLength = 8

	// Refresh token generation
	RefreshTokenBytes = 32

	// OTP settings (for future use)
	OTPLength     = 6
	OTPExpiration = 5 * time.Minute

	// Rate limiting (for future Redis implementation)
	MaxLoginAttemptsPerHour = 5
	LoginAttemptWindow      = 1 * time.Hour

	// Database
	DefaultMaxConns = 25
	DefaultMinConns = 5

	// Server
	DefaultPort = "8080"
)

// Error messages
const (
	ErrPhoneAlreadyRegistered = "phone number already registered"
	ErrEmailAlreadyRegistered = "email already registered"
	ErrInvalidCredentials     = "invalid phone number or password"
	ErrAccountSuspended       = "account is suspended"
	ErrInvalidRefreshToken    = "invalid refresh token"
	ErrTokenRevoked           = "token has been revoked"
	ErrTokenExpired           = "token has expired"
	ErrUserNotFound           = "user not found"
	ErrFailedToHashPassword   = "failed to hash password"
	ErrFailedToCreateUser     = "failed to create user"
	ErrFailedToGenerateToken  = "failed to generate token"
)

// Revoke reasons
const (
	RevokeReasonLogout         = "LOGOUT"
	RevokeReasonPasswordChange = "PASSWORD_CHANGE"
	RevokeReasonSecurity       = "SECURITY"
	RevokeReasonExpired        = "EXPIRED"
)
