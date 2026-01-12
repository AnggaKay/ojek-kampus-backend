package entity

import "time"

// PassengerProfile represents the passenger_profiles table
type PassengerProfile struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ProfilePicture     *string   `json:"profile_picture" db:"profile_picture"`
	FCMToken           *string   `json:"fcm_token" db:"fcm_token"`
	TotalOrders        int       `json:"total_orders" db:"total_orders"`
	TotalCancellations int       `json:"total_cancellations" db:"total_cancellations"`
	RatingAvg          float64   `json:"rating_avg" db:"rating_avg"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
