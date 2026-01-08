package middleware

import (
	"net/http"
	"strings"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	jwtPkg "github.com/AnggaKay/ojek-kampus-backend/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuth validates JWT token from Authorization header
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Missing authorization token"))
			}

			// Check Bearer format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Invalid authorization format"))
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtPkg.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Invalid or expired token"))
			}

			// Store claims in context
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			c.Set("user_type", claims.UserType)

			return next(c)
		}
	}
}

// RoleGuard restricts access to specific roles
func RoleGuard(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userType := c.Get("user_type")
			if userType == nil {
				return c.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", "Access denied"))
			}

			userTypeStr := userType.(string)
			allowed := false
			for _, role := range allowedRoles {
				if userTypeStr == role {
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", "Insufficient permissions"))
			}

			return next(c)
		}
	}
}
