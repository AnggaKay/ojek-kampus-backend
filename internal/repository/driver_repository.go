package repository

import (
	"context"
	"fmt"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriverRepository interface {
	Create(ctx context.Context, profile *entity.DriverProfile) error
	FindByID(ctx context.Context, id int) (*entity.DriverProfile, error)
	FindByUserID(ctx context.Context, userID int) (*entity.DriverProfile, error)
	ExistsByVehiclePlate(ctx context.Context, vehiclePlate string) (bool, error)
	Update(ctx context.Context, profile *entity.DriverProfile) error
	UpdateVerificationStatus(ctx context.Context, profileID int, isVerified bool, notes, reason *string, verifiedBy *int) error
}

type driverRepository struct {
	db *pgxpool.Pool
}

func NewDriverRepository(db *pgxpool.Pool) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) Create(ctx context.Context, profile *entity.DriverProfile) error {
	query := `
		INSERT INTO driver_profiles (
			user_id, vehicle_type, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color,
			ktp_photo, sim_photo, stnk_photo, ktm_photo, is_verified, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		profile.UserID,
		profile.VehicleType,
		profile.VehiclePlate,
		profile.VehicleBrand,
		profile.VehicleModel,
		profile.VehicleColor,
		profile.KTPPhoto,
		profile.SIMPhoto,
		profile.STNKPhoto,
		profile.KTMPhoto,
		profile.IsVerified,
		profile.IsActive,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
}

func (r *driverRepository) FindByID(ctx context.Context, id int) (*entity.DriverProfile, error) {
	query := `
		SELECT id, user_id, profile_picture, fcm_token, vehicle_type, vehicle_plate,
		       vehicle_brand, vehicle_model, vehicle_color,
		       ktp_photo, sim_photo, stnk_photo, ktm_photo,
		       is_verified, verification_notes, verified_by, verified_at, rejection_reason,
		       is_active, current_lat, current_long, last_location_update,
		       total_completed_orders, total_cancelled_orders, rating_avg,
		       created_at, updated_at
		FROM driver_profiles WHERE id = $1
	`

	var profile entity.DriverProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.ProfilePicture,
		&profile.FCMToken,
		&profile.VehicleType,
		&profile.VehiclePlate,
		&profile.VehicleBrand,
		&profile.VehicleModel,
		&profile.VehicleColor,
		&profile.KTPPhoto,
		&profile.SIMPhoto,
		&profile.STNKPhoto,
		&profile.KTMPhoto,
		&profile.IsVerified,
		&profile.VerificationNotes,
		&profile.VerifiedBy,
		&profile.VerifiedAt,
		&profile.RejectionReason,
		&profile.IsActive,
		&profile.CurrentLat,
		&profile.CurrentLong,
		&profile.LastLocationUpdate,
		&profile.TotalCompletedOrders,
		&profile.TotalCancelledOrders,
		&profile.RatingAvg,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("driver profile not found")
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *driverRepository) FindByUserID(ctx context.Context, userID int) (*entity.DriverProfile, error) {
	query := `
		SELECT id, user_id, profile_picture, fcm_token, vehicle_type, vehicle_plate,
		       vehicle_brand, vehicle_model, vehicle_color,
		       ktp_photo, sim_photo, stnk_photo, ktm_photo,
		       is_verified, verification_notes, verified_by, verified_at, rejection_reason,
		       is_active, current_lat, current_long, last_location_update,
		       total_completed_orders, total_cancelled_orders, rating_avg,
		       created_at, updated_at
		FROM driver_profiles WHERE user_id = $1
	`

	var profile entity.DriverProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.ProfilePicture,
		&profile.FCMToken,
		&profile.VehicleType,
		&profile.VehiclePlate,
		&profile.VehicleBrand,
		&profile.VehicleModel,
		&profile.VehicleColor,
		&profile.KTPPhoto,
		&profile.SIMPhoto,
		&profile.STNKPhoto,
		&profile.KTMPhoto,
		&profile.IsVerified,
		&profile.VerificationNotes,
		&profile.VerifiedBy,
		&profile.VerifiedAt,
		&profile.RejectionReason,
		&profile.IsActive,
		&profile.CurrentLat,
		&profile.CurrentLong,
		&profile.LastLocationUpdate,
		&profile.TotalCompletedOrders,
		&profile.TotalCancelledOrders,
		&profile.RatingAvg,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("driver profile not found")
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *driverRepository) Update(ctx context.Context, profile *entity.DriverProfile) error {
	query := `
		UPDATE driver_profiles
		SET profile_picture = $1, vehicle_brand = $2, vehicle_model = $3, vehicle_color = $4,
		    ktp_photo = $5, sim_photo = $6, stnk_photo = $7, ktm_photo = $8,
		    updated_at = NOW()
		WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query,
		profile.ProfilePicture,
		profile.VehicleBrand,
		profile.VehicleModel,
		profile.VehicleColor,
		profile.KTPPhoto,
		profile.SIMPhoto,
		profile.STNKPhoto,
		profile.KTMPhoto,
		profile.ID,
	)
	return err
}

// ExistsByVehiclePlate checks if vehicle plate already exists
func (r *driverRepository) ExistsByVehiclePlate(ctx context.Context, plate string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM driver_profiles WHERE vehicle_plate = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, plate).Scan(&exists)
	return exists, err
}

func (r *driverRepository) UpdateVerificationStatus(ctx context.Context, profileID int, isVerified bool, notes, reason *string, verifiedBy *int) error {
	query := `
		UPDATE driver_profiles
		SET is_verified = $1, verification_notes = $2, rejection_reason = $3, 
		    verified_by = $4, verified_at = CASE WHEN $1 = TRUE THEN NOW() ELSE NULL END,
		    updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, isVerified, notes, reason, verifiedBy, profileID)
	return err
}
