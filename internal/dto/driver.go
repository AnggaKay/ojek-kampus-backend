package dto

// ============================================================================
// Driver Request DTOs
// ============================================================================

// RegisterDriverRequest represents driver registration request (multipart form)
type RegisterDriverRequest struct {
	PhoneNumber  string `form:"phone_number" validate:"required"`
	Password     string `form:"password" validate:"required,min=8"`
	FullName     string `form:"full_name" validate:"required"`
	Email        string `form:"email" validate:"omitempty,email"`
	VehiclePlate string `form:"vehicle_plate" validate:"required"`
	VehicleBrand string `form:"vehicle_brand"`
	VehicleModel string `form:"vehicle_model"`
	VehicleColor string `form:"vehicle_color"`
	// Files are handled separately via multipart form
}

// UpdateDriverProfileRequest represents driver profile update request
type UpdateDriverProfileRequest struct {
	VehiclePlate *string `json:"vehicle_plate,omitempty"`
	VehicleBrand *string `json:"vehicle_brand,omitempty"`
	VehicleModel *string `json:"vehicle_model,omitempty"`
	VehicleColor *string `json:"vehicle_color,omitempty"`
}

// ============================================================================
// Driver Response DTOs
// ============================================================================

// DriverProfileResponse represents driver profile data
type DriverProfileResponse struct {
	ID                   int        `json:"id"`
	UserID               int        `json:"user_id"`
	VehicleType          string     `json:"vehicle_type"`
	VehiclePlate         string     `json:"vehicle_plate"`
	VehicleBrand         *string    `json:"vehicle_brand,omitempty"`
	VehicleModel         *string    `json:"vehicle_model,omitempty"`
	VehicleColor         *string    `json:"vehicle_color,omitempty"`
	IsVerified           bool       `json:"is_verified"`
	VerificationStatus   string     `json:"verification_status"`
	RejectionReason      *string    `json:"rejection_reason,omitempty"`
	IsActive             bool       `json:"is_active"`
	TotalCompletedOrders int        `json:"total_completed_orders"`
	RatingAvg            float64    `json:"rating_avg"`
	Documents            *Documents `json:"documents,omitempty"`
}

// DriverAuthResponse represents driver authentication response (for registration)
type DriverAuthResponse struct {
	User          *UserResponse          `json:"user"`
	DriverProfile *DriverProfileResponse `json:"driver_profile"`
	AccessToken   string                 `json:"access_token"`
	RefreshToken  string                 `json:"refresh_token"`
	ExpiresIn     int                    `json:"expires_in"`
}
