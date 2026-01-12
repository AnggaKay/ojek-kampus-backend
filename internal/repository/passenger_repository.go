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
			profile_picture,
			fcm_token
		)
		VALUES ($1, $2, $3)
		RETURNING id, total_orders, total_cancellations, rating_avg, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		profile.UserID,
		profile.ProfilePicture,
		profile.FCMToken,
	).Scan(
		&profile.ID,
		&profile.TotalOrders,
		&profile.TotalCancellations,
		&profile.RatingAvg,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

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
			profile_picture,
			fcm_token,
			total_orders,
			total_cancellations,
			rating_avg,
			created_at,
			updated_at
		FROM passenger_profiles
		WHERE id = $1
	`

	profile := &entity.PassengerProfile{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.ProfilePicture,
		&profile.FCMToken,
		&profile.TotalOrders,
		&profile.TotalCancellations,
		&profile.RatingAvg,
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
			profile_picture,
			fcm_token,
			total_orders,
			total_cancellations,
			rating_avg,
			created_at,
			updated_at
		FROM passenger_profiles
		WHERE user_id = $1
	`

	profile := &entity.PassengerProfile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.ProfilePicture,
		&profile.FCMToken,
		&profile.TotalOrders,
		&profile.TotalCancellations,
		&profile.RatingAvg,
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
			profile_picture = $2,
			fcm_token = $3,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		profile.ID,
		profile.ProfilePicture,
		profile.FCMToken,
	).Scan(&profile.UpdatedAt)

	if err != nil {
		logger.Log.Error().Err(err).Int("passenger_id", profile.ID).Msg("Failed to update passenger profile")
		return err
	}

	logger.Log.Info().Int("passenger_id", profile.ID).Msg("Passenger profile updated")
	return nil
}
