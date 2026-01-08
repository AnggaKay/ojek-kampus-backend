package entity

import "time"

type RefreshToken struct {
	ID           int        `json:"id" db:"id"`
	UserID       int        `json:"user_id" db:"user_id"`
	UserType     string     `json:"user_type" db:"user_type"` // PASSENGER, DRIVER, ADMIN
	TokenHash    string     `json:"-" db:"token_hash"`
	DeviceInfo   *string    `json:"device_info,omitempty" db:"device_info"`
	DeviceName   *string    `json:"device_name,omitempty" db:"device_name"`
	IPAddress    *string    `json:"ip_address,omitempty" db:"ip_address"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
	IsRevoked    bool       `json:"is_revoked" db:"is_revoked"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokeReason *string    `json:"revoke_reason,omitempty" db:"revoke_reason"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt   time.Time  `json:"last_used_at" db:"last_used_at"`
}
