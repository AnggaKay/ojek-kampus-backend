package entity

import "time"

// DriverProfile represents driver-specific profile data
type DriverProfile struct {
	ID                   int        `db:"id"`
	UserID               int        `db:"user_id"`
	ProfilePicture       *string    `db:"profile_picture"`
	FCMToken             *string    `db:"fcm_token"`
	VehicleType          string     `db:"vehicle_type"`
	VehiclePlate         string     `db:"vehicle_plate"`
	VehicleBrand         *string    `db:"vehicle_brand"`
	VehicleModel         *string    `db:"vehicle_model"`
	VehicleColor         *string    `db:"vehicle_color"`
	KTPPhoto             *string    `db:"ktp_photo"`
	SIMPhoto             *string    `db:"sim_photo"`
	STNKPhoto            *string    `db:"stnk_photo"`
	KTMPhoto             *string    `db:"ktm_photo"`
	IsVerified           bool       `db:"is_verified"`
	VerificationNotes    *string    `db:"verification_notes"`
	VerifiedBy           *int       `db:"verified_by"`
	VerifiedAt           *time.Time `db:"verified_at"`
	RejectionReason      *string    `db:"rejection_reason"`
	IsActive             bool       `db:"is_active"`
	CurrentLat           *float64   `db:"current_lat"`
	CurrentLong          *float64   `db:"current_long"`
	LastLocationUpdate   *time.Time `db:"last_location_update"`
	TotalCompletedOrders int        `db:"total_completed_orders"`
	TotalCancelledOrders int        `db:"total_cancelled_orders"`
	RatingAvg            float64    `db:"rating_avg"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

// VehicleInfo represents vehicle information for registration
type VehicleInfo struct {
	Type  string
	Plate string
	Brand string
	Model string
	Color string
}

// DriverDocuments represents uploaded document paths
type DriverDocuments struct {
	KTPPhoto  string
	SIMPhoto  string
	STNKPhoto string
	KTMPhoto  string
}
