package handler

import (
	"net/http"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.SuccessResponse("Service is healthy", map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "ojek-kampus-backend",
	}))
}
