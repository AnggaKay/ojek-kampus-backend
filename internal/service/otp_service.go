package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/entity"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/whatsapp"
)

const (
	OTPExpireMinutes  = 5
	OTPResendCooldown = 60 // seconds
	MaxVerifyAttempts = 3
)

// OTPService handles OTP business logic
type OTPService interface {
	SendOTP(ctx context.Context, req dto.SendOTPRequest, ipAddress, userAgent string) (*dto.SendOTPResponse, error)
	VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error)
	ResendOTP(ctx context.Context, req dto.ResendOTPRequest, ipAddress, userAgent string) (*dto.SendOTPResponse, error)
}

type otpService struct {
	otpRepo        repository.OTPRepository
	whatsappClient *whatsapp.WhatsAppClient
}

// NewOTPService creates a new OTP service
func NewOTPService(otpRepo repository.OTPRepository, whatsappClient *whatsapp.WhatsAppClient) OTPService {
	return &otpService{
		otpRepo:        otpRepo,
		whatsappClient: whatsappClient,
	}
}

// SendOTP generates and sends OTP code to user's WhatsApp
func (s *otpService) SendOTP(ctx context.Context, req dto.SendOTPRequest, ipAddress, userAgent string) (*dto.SendOTPResponse, error) {
	logger.Log.Info().
		Str("phone", req.PhoneNumber).
		Str("purpose", string(req.Purpose)).
		Msg("Sending OTP")

	// Check if there's a recent OTP (rate limiting)
	existingOTP, err := s.otpRepo.FindLatestByPhoneAndPurpose(ctx, req.PhoneNumber, req.Purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing OTP: %w", err)
	}

	if existingOTP != nil && !existingOTP.IsExpired() {
		// Check cooldown period
		timeSinceCreation := time.Since(existingOTP.CreatedAt).Seconds()
		if timeSinceCreation < OTPResendCooldown {
			remainingTime := int(OTPResendCooldown - timeSinceCreation)
			return nil, fmt.Errorf("please wait %d seconds before requesting new OTP", remainingTime)
		}
	}

	// Invalidate old OTPs
	if err := s.otpRepo.InvalidateOldOTPs(ctx, req.PhoneNumber, req.Purpose); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to invalidate old OTPs")
	}

	// Generate 6-digit OTP
	otpCode, err := s.generateOTP()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create OTP record
	now := time.Now()
	otp := &entity.OTPCode{
		PhoneNumber: req.PhoneNumber,
		OTPCode:     otpCode,
		Purpose:     req.Purpose,
		ExpiresAt:   now.Add(OTPExpireMinutes * time.Minute),
		IsUsed:      false,
		Attempts:    0,
		IPAddress:   &ipAddress,
		UserAgent:   &userAgent,
		CreatedAt:   now,
	}

	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP via WhatsApp
	if err := s.whatsappClient.SendOTP(req.PhoneNumber, otpCode); err != nil {
		logger.Log.Error().
			Err(err).
			Str("phone", req.PhoneNumber).
			Msg("Failed to send WhatsApp message")
		return nil, fmt.Errorf("failed to send OTP: %w", err)
	}

	return &dto.SendOTPResponse{
		PhoneNumber: req.PhoneNumber,
		ExpiresIn:   OTPExpireMinutes * 60, // in seconds
		Message:     "Kode OTP telah dikirim ke WhatsApp Anda",
	}, nil
}

// VerifyOTP verifies the OTP code provided by user
func (s *otpService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	logger.Log.Info().
		Str("phone", req.PhoneNumber).
		Msg("Verifying OTP")

	// Find OTP by phone and code
	otp, err := s.otpRepo.FindByPhoneAndCode(ctx, req.PhoneNumber, req.OTPCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}

	if otp == nil {
		return &dto.VerifyOTPResponse{
			PhoneNumber: req.PhoneNumber,
			Verified:    false,
			Message:     "Kode OTP tidak valid",
		}, nil
	}

	// Increment attempts
	if err := s.otpRepo.IncrementAttempts(ctx, otp.ID); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to increment attempts")
	}

	// Check if already used
	if otp.IsUsed {
		return &dto.VerifyOTPResponse{
			PhoneNumber: req.PhoneNumber,
			Verified:    false,
			Message:     "Kode OTP sudah pernah digunakan",
		}, nil
	}

	// Check if expired
	if otp.IsExpired() {
		return &dto.VerifyOTPResponse{
			PhoneNumber: req.PhoneNumber,
			Verified:    false,
			Message:     "Kode OTP telah kadaluarsa",
		}, nil
	}

	// Check attempts limit
	if otp.Attempts >= MaxVerifyAttempts {
		return &dto.VerifyOTPResponse{
			PhoneNumber: req.PhoneNumber,
			Verified:    false,
			Message:     "Terlalu banyak percobaan. Silakan minta kode baru",
		}, nil
	}

	// Mark as used
	if err := s.otpRepo.MarkAsUsed(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	logger.Log.Info().
		Str("phone", req.PhoneNumber).
		Msg("OTP verified successfully")

	return &dto.VerifyOTPResponse{
		PhoneNumber: req.PhoneNumber,
		Verified:    true,
		Message:     "Verifikasi berhasil",
	}, nil
}

// ResendOTP resends OTP code (same as SendOTP with cooldown check)
func (s *otpService) ResendOTP(ctx context.Context, req dto.ResendOTPRequest, ipAddress, userAgent string) (*dto.SendOTPResponse, error) {
	return s.SendOTP(ctx, dto.SendOTPRequest{
		PhoneNumber: req.PhoneNumber,
		Purpose:     req.Purpose,
	}, ipAddress, userAgent)
}

// generateOTP generates a secure 6-digit OTP code
func (s *otpService) generateOTP() (string, error) {
	// Generate random number between 100000 and 999999
	min := big.NewInt(100000)
	max := big.NewInt(999999)
	diff := new(big.Int).Sub(max, min)

	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return "", err
	}

	otp := new(big.Int).Add(min, n)
	return fmt.Sprintf("%06d", otp.Int64()), nil
}
