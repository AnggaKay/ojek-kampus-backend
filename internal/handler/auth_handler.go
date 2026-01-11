package handler

import (
	"net/http"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/middleware"
	"github.com/AnggaKay/ojek-kampus-backend/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterPassenger godoc
// @Summary Register new passenger
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterPassengerRequest true "Registration data"
// @Success 201 {object} dto.Response
// @Router /api/auth/register/passenger [post]
func (h *AuthHandler) RegisterPassenger(c echo.Context) error {
	var req dto.RegisterPassengerRequest

	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	result, err := h.authService.RegisterPassenger(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("REGISTRATION_FAILED", err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse("Registration successful", result))
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	// Add device info from User-Agent
	if req.DeviceInfo == "" {
		req.DeviceInfo = c.Request().UserAgent()
	}

	result, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("LOGIN_FAILED", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Login successful", result))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.Response
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest

	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	result, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("REFRESH_FAILED", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Token refreshed", result))
}

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.RefreshTokenRequest

	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	if err := h.authService.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("LOGOUT_FAILED", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Logout successful", nil))
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Router /api/auth/me [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Invalid user context"))
	}

	userRole := c.Get("user_role")

	return c.JSON(http.StatusOK, dto.SuccessResponse("Profile retrieved", map[string]interface{}{
		"user_id": userID,
		"role":    userRole,
	}))
}
