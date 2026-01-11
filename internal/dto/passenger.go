package dto

// ============================================================================
// Passenger Request DTOs
// ============================================================================

// RegisterPassengerRequest represents passenger registration request
type RegisterPassengerRequest struct {
	PhoneNumber string  `json:"phone_number" validate:"required,min=10,max=15"`
	Password    string  `json:"password" validate:"required,min=8"`
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UpdatePassengerProfileRequest represents passenger profile update request
type UpdatePassengerProfileRequest struct {
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	HomeAddress           *string `json:"home_address,omitempty"`
}

// ============================================================================
// Passenger Response DTOs
// ============================================================================

// PassengerProfileResponse represents passenger profile data
type PassengerProfileResponse struct {
	ID                    int     `json:"id"`
	UserID                int     `json:"user_id"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	HomeAddress           *string `json:"home_address,omitempty"`
	TotalCompletedOrders  int     `json:"total_completed_orders"`
}
