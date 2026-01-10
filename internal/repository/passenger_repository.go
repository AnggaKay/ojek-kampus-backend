package repository

import (
	"context"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassengerRepository interface {
	Create(ctx context.Context, profile *entity.PassengerProfile) error
	FindByID(ctx context.Context, id int) (*entity.PassengerProfile, error)
	FindByUserID(ctx context.Context, userID int) (*entity.PassengerProfile, error)
	Update(ctx context.Context, profile *entity.PassengerProfile) error
}

type passengerRepository struct {
	db *pgxpool.Pool
}

func NewPassengerRepository(db *pgxpool.Pool) PassengerRepository {
	return &passengerRepository{db: db}
}

func (r *passengerRepository) Create(ctx context.Context, profile *entity.PassengerProfile) error {
	query := `
		INSERT INTO passenger_profiles (
			user_id, 
			emergency_contact_name, 
			emergency_contact_phone, 
			home_address
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		profile.UserID,
		profile.EmergencyContactName,
		profile.EmergencyContactPhone,
		profile.HomeAddress,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", profile.UserID).Msg("Failed to create passenger profile")
		return err
	}

	logger.Log.Info().Int("passenger_id", profile.ID).Int("user_id", profile.UserID).Msg("Passenger profile created")
	return nil
}

func (r *passengerRepository) FindByID(ctx context.Context, id int) (*entity.PassengerProfile, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			emergency_contact_name, 
			emergency_contact_phone, 
			home_address,
			total_completed_orders,
			created_at, 
			updated_at
		FROM passenger_profiles
		WHERE id = $1
	`

	profile := &entity.PassengerProfile{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.EmergencyContactName,
		&profile.EmergencyContactPhone,
		&profile.HomeAddress,
		&profile.TotalCompletedOrders,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		logger.Log.Error().Err(err).Int("passenger_id", id).Msg("Failed to find passenger profile by ID")
		return nil, err
	}

	return profile, nil
}

func (r *passengerRepository) FindByUserID(ctx context.Context, userID int) (*entity.PassengerProfile, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			emergency_contact_name, 
			emergency_contact_phone, 
			home_address,
			total_completed_orders,
			created_at, 
			updated_at
		FROM passenger_profiles
		WHERE user_id = $1
	`

	profile := &entity.PassengerProfile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.EmergencyContactName,
		&profile.EmergencyContactPhone,
		&profile.HomeAddress,
		&profile.TotalCompletedOrders,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("Failed to find passenger profile by user ID")
		return nil, err
	}

	return profile, nil
}

func (r *passengerRepository) Update(ctx context.Context, profile *entity.PassengerProfile) error {
	query := `
		UPDATE passenger_profiles
		SET 
			emergency_contact_name = $2,
			emergency_contact_phone = $3,
			home_address = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		profile.ID,
		profile.EmergencyContactName,
		profile.EmergencyContactPhone,
		profile.HomeAddress,
	).Scan(&profile.UpdatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Int("passenger_id", profile.ID).Msg("Failed to update passenger profile")
		return err
	}

	logger.Log.Info().Int("passenger_id", profile.ID).Msg("Passenger profile updated")
	return nil
}
