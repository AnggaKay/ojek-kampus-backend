package middleware

import (
	"net/http"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// ValidateRequest binds and validates request body
func ValidateRequest(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if err := c.Validate(req); err != nil {
		validationErrors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				field := fe.Field()
				switch fe.Tag() {
				case "required":
					validationErrors[field] = field + " is required"
				case "email":
					validationErrors[field] = field + " must be a valid email"
				case "min":
					validationErrors[field] = field + " is too short"
				case "max":
					validationErrors[field] = field + " is too long"
				default:
					validationErrors[field] = field + " is invalid"
				}
			}
		}

		return c.JSON(http.StatusBadRequest, dto.ValidationErrorResponse(validationErrors))
	}

	return nil
}
