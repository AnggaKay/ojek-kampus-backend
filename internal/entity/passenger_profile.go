package entity

import "time"

// PassengerProfile represents the passenger_profiles table
type PassengerProfile struct {
	ID                    int       `json:"id" db:"id"`
	UserID                int       `json:"user_id" db:"user_id"`
	EmergencyContactName  *string   `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone *string   `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	HomeAddress           *string   `json:"home_address" db:"home_address"`
	TotalCompletedOrders  int       `json:"total_completed_orders" db:"total_completed_orders"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}
