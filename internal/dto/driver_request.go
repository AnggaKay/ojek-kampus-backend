package dto

// RegisterDriverRequest represents driver registration request
type RegisterDriverRequest struct {
	PhoneNumber  string `form:"phone_number" validate:"required"`
	Password     string `form:"password" validate:"required,min=8"`
	FullName     string `form:"full_name" validate:"required"`
	Email        string `form:"email" validate:"omitempty,email"`
	VehiclePlate string `form:"vehicle_plate" validate:"required"`
	VehicleBrand string `form:"vehicle_brand"`
	VehicleModel string `form:"vehicle_model"`
	VehicleColor string `form:"vehicle_color"`
	// Files will be handled separately via multipart form
}

// DriverProfileResponse represents driver profile data
type DriverProfileResponse struct {
	ID                   int        `json:"id"`
	UserID               int        `json:"user_id"`
	VehicleType          string     `json:"vehicle_type"`
	VehiclePlate         string     `json:"vehicle_plate"`
	VehicleBrand         *string    `json:"vehicle_brand"`
	VehicleModel         *string    `json:"vehicle_model"`
	VehicleColor         *string    `json:"vehicle_color"`
	IsVerified           bool       `json:"is_verified"`
	VerificationStatus   string     `json:"verification_status"`
	RejectionReason      *string    `json:"rejection_reason,omitempty"`
	IsActive             bool       `json:"is_active"`
	TotalCompletedOrders int        `json:"total_completed_orders"`
	RatingAvg            float64    `json:"rating_avg"`
	Documents            *Documents `json:"documents,omitempty"`
}

// Documents represents document upload status
type Documents struct {
	KTPUploaded  bool `json:"ktp_uploaded"`
	SIMUploaded  bool `json:"sim_uploaded"`
	STNKUploaded bool `json:"stnk_uploaded"`
	KTMUploaded  bool `json:"ktm_uploaded"`
}

// DriverAuthResponse represents driver authentication response
type DriverAuthResponse struct {
	User          *UserResponse          `json:"user"`
	DriverProfile *DriverProfileResponse `json:"driver_profile"`
	AccessToken   string                 `json:"access_token"`
	RefreshToken  string                 `json:"refresh_token"`
	ExpiresIn     int                    `json:"expires_in"`
}
