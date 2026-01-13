package handler

import (
	"net/http"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/service"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

// OTPHandler handles OTP-related HTTP requests
type OTPHandler struct {
	otpService service.OTPService
}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler(otpService service.OTPService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

// SendOTP godoc
// @Summary Send OTP code
// @Description Send OTP code to user's WhatsApp number
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.SendOTPRequest true "Send OTP Request"
// @Success 200 {object} dto.Response{data=dto.SendOTPResponse}
// @Failure 400 {object} dto.Response
// @Failure 429 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/auth/send-otp [post]
func (h *OTPHandler) SendOTP(c echo.Context) error {
	var req dto.SendOTPRequest
	if err := c.Bind(&req); err != nil {
		logger.Log.Warn().Err(err).Msg("Invalid request body")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Get IP address and user agent
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Send OTP
	response, err := h.otpService.SendOTP(c.Request().Context(), req, ipAddress, userAgent)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to send OTP")

		// Check for rate limit error
		if contains(err.Error(), "please wait") {
			return c.JSON(http.StatusTooManyRequests, dto.Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal mengirim kode OTP. Silakan coba lagi.",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "OTP sent successfully",
		Data:    response,
	})
}

// VerifyOTP godoc
// @Summary Verify OTP code
// @Description Verify OTP code entered by user
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} dto.Response{data=dto.VerifyOTPResponse}
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/auth/verify-otp [post]
func (h *OTPHandler) VerifyOTP(c echo.Context) error {
	var req dto.VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		logger.Log.Warn().Err(err).Msg("Invalid request body")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Verify OTP
	response, err := h.otpService.VerifyOTP(c.Request().Context(), req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to verify OTP")
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal memverifikasi kode OTP. Silakan coba lagi.",
			Data:    nil,
		})
	}

	// If verification failed, return 400
	if !response.Verified {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: response.Message,
			Data:    response,
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "OTP verified successfully",
		Data:    response,
	})
}

// ResendOTP godoc
// @Summary Resend OTP code
// @Description Resend OTP code to user's WhatsApp number
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.ResendOTPRequest true "Resend OTP Request"
// @Success 200 {object} dto.Response{data=dto.SendOTPResponse}
// @Failure 400 {object} dto.Response
// @Failure 429 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/auth/resend-otp [post]
func (h *OTPHandler) ResendOTP(c echo.Context) error {
	var req dto.ResendOTPRequest
	if err := c.Bind(&req); err != nil {
		logger.Log.Warn().Err(err).Msg("Invalid request body")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Get IP address and user agent
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Resend OTP
	response, err := h.otpService.ResendOTP(c.Request().Context(), req, ipAddress, userAgent)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to resend OTP")

		// Check for rate limit error
		if contains(err.Error(), "please wait") {
			return c.JSON(http.StatusTooManyRequests, dto.Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal mengirim ulang kode OTP. Silakan coba lagi.",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "OTP resent successfully",
		Data:    response,
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
