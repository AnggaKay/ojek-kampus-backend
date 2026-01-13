package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OTPRepository handles OTP database operations
type OTPRepository interface {
	Create(ctx context.Context, otp *entity.OTPCode) error
	FindLatestByPhoneAndPurpose(ctx context.Context, phoneNumber string, purpose entity.OTPPurpose) (*entity.OTPCode, error)
	FindByPhoneAndCode(ctx context.Context, phoneNumber, otpCode string) (*entity.OTPCode, error)
	MarkAsUsed(ctx context.Context, id int) error
	IncrementAttempts(ctx context.Context, id int) error
	InvalidateOldOTPs(ctx context.Context, phoneNumber string, purpose entity.OTPPurpose) error
}

type otpRepository struct {
	db *pgxpool.Pool
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *pgxpool.Pool) OTPRepository {
	return &otpRepository{db: db}
}

// Create inserts a new OTP code
func (r *otpRepository) Create(ctx context.Context, otp *entity.OTPCode) error {
	query := `
		INSERT INTO otp_codes (
			phone_number, otp_code, purpose, expires_at, 
			ip_address, user_agent, attempts, is_used, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, 0, false, $7
		) RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		otp.PhoneNumber,
		otp.OTPCode,
		otp.Purpose,
		otp.ExpiresAt,
		otp.IPAddress,
		otp.UserAgent,
		otp.CreatedAt,
	).Scan(&otp.ID)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("phone", otp.PhoneNumber).
			Msg("Failed to create OTP")
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	logger.Log.Info().
		Int("id", otp.ID).
		Str("phone", otp.PhoneNumber).
		Str("purpose", string(otp.Purpose)).
		Msg("OTP created successfully")

	return nil
}

// FindLatestByPhoneAndPurpose finds the latest OTP for a phone number and purpose
func (r *otpRepository) FindLatestByPhoneAndPurpose(ctx context.Context, phoneNumber string, purpose entity.OTPPurpose) (*entity.OTPCode, error) {
	query := `
		SELECT id, phone_number, otp_code, purpose, expires_at, 
		       is_used, used_at, attempts, ip_address, user_agent, created_at
		FROM otp_codes
		WHERE phone_number = $1 AND purpose = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var otp entity.OTPCode
	err := r.db.QueryRow(ctx, query, phoneNumber, purpose).Scan(
		&otp.ID,
		&otp.PhoneNumber,
		&otp.OTPCode,
		&otp.Purpose,
		&otp.ExpiresAt,
		&otp.IsUsed,
		&otp.UsedAt,
		&otp.Attempts,
		&otp.IPAddress,
		&otp.UserAgent,
		&otp.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		logger.Log.Error().
			Err(err).
			Str("phone", phoneNumber).
			Str("purpose", string(purpose)).
			Msg("Failed to find OTP")
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}

	return &otp, nil
}

// FindByPhoneAndCode finds OTP by phone number and code (for verification)
func (r *otpRepository) FindByPhoneAndCode(ctx context.Context, phoneNumber, otpCode string) (*entity.OTPCode, error) {
	query := `
		SELECT id, phone_number, otp_code, purpose, expires_at, 
		       is_used, used_at, attempts, ip_address, user_agent, created_at
		FROM otp_codes
		WHERE phone_number = $1 AND otp_code = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var otp entity.OTPCode
	err := r.db.QueryRow(ctx, query, phoneNumber, otpCode).Scan(
		&otp.ID,
		&otp.PhoneNumber,
		&otp.OTPCode,
		&otp.Purpose,
		&otp.ExpiresAt,
		&otp.IsUsed,
		&otp.UsedAt,
		&otp.Attempts,
		&otp.IPAddress,
		&otp.UserAgent,
		&otp.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		logger.Log.Error().
			Err(err).
			Str("phone", phoneNumber).
			Msg("Failed to find OTP by code")
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}

	return &otp, nil
}

// MarkAsUsed marks an OTP as used
func (r *otpRepository) MarkAsUsed(ctx context.Context, id int) error {
	query := `
		UPDATE otp_codes
		SET is_used = true, used_at = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Int("id", id).
			Msg("Failed to mark OTP as used")
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("OTP not found")
	}

	logger.Log.Info().Int("id", id).Msg("OTP marked as used")
	return nil
}

// IncrementAttempts increments the verification attempts counter
func (r *otpRepository) IncrementAttempts(ctx context.Context, id int) error {
	query := `
		UPDATE otp_codes
		SET attempts = attempts + 1
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Int("id", id).
			Msg("Failed to increment OTP attempts")
		return fmt.Errorf("failed to increment attempts: %w", err)
	}

	return nil
}

// InvalidateOldOTPs marks all old OTPs as used (when generating new OTP)
func (r *otpRepository) InvalidateOldOTPs(ctx context.Context, phoneNumber string, purpose entity.OTPPurpose) error {
	query := `
		UPDATE otp_codes
		SET is_used = true
		WHERE phone_number = $1 AND purpose = $2 AND is_used = false
	`

	_, err := r.db.Exec(ctx, query, phoneNumber, purpose)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("phone", phoneNumber).
			Str("purpose", string(purpose)).
			Msg("Failed to invalidate old OTPs")
		return fmt.Errorf("failed to invalidate old OTPs: %w", err)
	}

	return nil
}
