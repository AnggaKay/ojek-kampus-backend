package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	UpdateLastUsed(ctx context.Context, id int) error
	Revoke(ctx context.Context, id int, reason string) error
	RevokeAllByUserID(ctx context.Context, userID int, userType string) error
	DeleteExpired(ctx context.Context) error
}

type refreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, user_type, token_hash, device_info, device_name, 
		                            ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, last_used_at
	`
	return r.db.QueryRow(ctx, query,
		token.UserID,
		token.UserType,
		token.TokenHash,
		token.DeviceInfo,
		token.DeviceName,
		token.IPAddress,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt, &token.LastUsedAt)
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, user_type, token_hash, device_info, device_name, ip_address,
		       expires_at, is_revoked, revoked_at, revoke_reason, created_at, last_used_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var token entity.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.UserType,
		&token.TokenHash,
		&token.DeviceInfo,
		&token.DeviceName,
		&token.IPAddress,
		&token.ExpiresAt,
		&token.IsRevoked,
		&token.RevokedAt,
		&token.RevokeReason,
		&token.CreatedAt,
		&token.LastUsedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("token not found")
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) UpdateLastUsed(ctx context.Context, id int) error {
	query := `UPDATE refresh_tokens SET last_used_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id int, reason string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = NOW(), revoke_reason = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, reason, id)
	return err
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID int, userType string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = NOW(), revoke_reason = 'LOGOUT_ALL'
		WHERE user_id = $1 AND user_type = $2 AND is_revoked = false
	`
	_, err := r.db.Exec(ctx, query, userID, userType)
	return err
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(ctx, query, time.Now())
	return err
}
