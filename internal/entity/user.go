package entity

import "time"

type UserRole string
type UserStatus string

const (
	RolePassenger UserRole = "PASSENGER"
	RoleDriver    UserRole = "DRIVER"
)

const (
	StatusActive              UserStatus = "ACTIVE"
	StatusSuspended           UserStatus = "SUSPENDED"
	StatusPendingVerification UserStatus = "PENDING_VERIFICATION"
	StatusRejected            UserStatus = "REJECTED"
)

type User struct {
	ID            int        `json:"id" db:"id"`
	PhoneNumber   string     `json:"phone_number" db:"phone_number"`
	PasswordHash  string     `json:"-" db:"password_hash"` // Hidden from JSON
	Email         *string    `json:"email,omitempty" db:"email"`
	FullName      string     `json:"full_name" db:"full_name"`
	Role          UserRole   `json:"role" db:"role"`
	Status        UserStatus `json:"status" db:"status"`
	PhoneVerified bool       `json:"phone_verified" db:"phone_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}
