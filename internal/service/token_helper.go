package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
)

// TokenHelper provides shared token operations
type TokenHelper struct {
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewTokenHelper(refreshTokenRepo repository.RefreshTokenRepository) *TokenHelper {
	return &TokenHelper{
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (h *TokenHelper) CreateRefreshToken(ctx context.Context, userID int, userType, deviceInfo string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, constants.RefreshTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate random bytes for refresh token")
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	tokenHash := HashToken(token)

	// Create token record
	refreshToken := &entity.RefreshToken{
		UserID:     userID,
		UserType:   userType,
		TokenHash:  tokenHash,
		DeviceInfo: &deviceInfo,
		ExpiresAt:  time.Now().Add(constants.RefreshTokenTTL),
		IsRevoked:  false,
	}

	if err := h.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		logger.Log.Error().Err(err).Int("user_id", userID).Msg("Failed to save refresh token")
		return "", err
	}

	return token, nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
