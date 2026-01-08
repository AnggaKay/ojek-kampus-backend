package repository

import (
	"context"
	"fmt"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int) (*entity.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateLastLogin(ctx context.Context, userID int) error
	UpdatePhoneVerified(ctx context.Context, userID int, verified bool) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (phone_number, password_hash, email, full_name, role, status, phone_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		user.PhoneNumber,
		user.PasswordHash,
		user.Email,
		user.FullName,
		user.Role,
		user.Status,
		user.PhoneVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
		SELECT id, phone_number, password_hash, email, full_name, role, status, 
		       phone_verified, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.PhoneVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	query := `
		SELECT id, phone_number, password_hash, email, full_name, role, status,
		       phone_verified, last_login_at, created_at, updated_at
		FROM users WHERE phone_number = $1
	`
	var user entity.User
	err := r.db.QueryRow(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.PhoneVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, phone_number, password_hash, email, full_name, role, status,
		       phone_verified, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.PhoneVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET email = $1, full_name = $2, status = $3, phone_verified = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		user.Email,
		user.FullName,
		user.Status,
		user.PhoneVerified,
		user.ID,
	)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *userRepository) UpdatePhoneVerified(ctx context.Context, userID int, verified bool) error {
	query := `UPDATE users SET phone_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, verified, userID)
	return err
}
