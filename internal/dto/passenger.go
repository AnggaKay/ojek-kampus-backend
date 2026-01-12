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
	ProfilePicture *string `json:"profile_picture,omitempty"`
	FCMToken       *string `json:"fcm_token,omitempty"`
}

// ============================================================================
// Passenger Response DTOs
// ============================================================================

// PassengerProfileResponse represents passenger profile data
type PassengerProfileResponse struct {
	ID                 int     `json:"id"`
	UserID             int     `json:"user_id"`
	ProfilePicture     *string `json:"profile_picture,omitempty"`
	FCMToken           *string `json:"fcm_token,omitempty"`
	TotalOrders        int     `json:"total_orders"`
	TotalCancellations int     `json:"total_cancellations"`
	RatingAvg          float64 `json:"rating_avg"`
}
