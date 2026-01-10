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

	// File Upload
	MaxFileSize        = 5 * 1024 * 1024  // 5MB
	MaxTotalUploadSize = 20 * 1024 * 1024 // 20MB (4 files x 5MB)
	UploadDirectory    = "uploads"
	AllowedImageTypes  = "image/jpeg,image/png"
	AllowedDocTypes    = "application/pdf"
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

	// Driver-specific errors
	ErrVehiclePlateExists = "vehicle plate already registered"
	ErrInvalidFileType    = "invalid file type"
	ErrFileTooLarge       = "file size exceeds maximum limit"
	ErrFailedToUploadFile = "failed to upload file"
	ErrDriverNotVerified  = "driver account is not verified yet"
	ErrDocumentNotFound   = "document not found"
	ErrUnauthorizedAccess = "unauthorized access to document"
)

// Revoke reasons
const (
	RevokeReasonLogout         = "LOGOUT"
	RevokeReasonPasswordChange = "PASSWORD_CHANGE"
	RevokeReasonSecurity       = "SECURITY"
	RevokeReasonExpired        = "EXPIRED"
)

// Vehicle types (Motor only - Ojek Kampus focuses on motorcycle transportation)
const (
	VehicleTypeMotor = "MOTOR"
)

// Document types
const (
	DocumentTypeKTP  = "KTP"
	DocumentTypeSIM  = "SIM"
	DocumentTypeSTNK = "STNK"
	DocumentTypeKTM  = "KTM"
)

// Driver verification status
const (
	VerificationStatusPending  = "PENDING_VERIFICATION"
	VerificationStatusVerified = "VERIFIED"
	VerificationStatusRejected = "REJECTED"
)
