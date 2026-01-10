package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	uploadDir string
}

func NewDocumentHandler(uploadDir string) *DocumentHandler {
	return &DocumentHandler{
		uploadDir: uploadDir,
	}
}

// GetDocument serves document files with authentication and authorization
// GET /api/documents/:type/:filename
func (h *DocumentHandler) GetDocument(c echo.Context) error {
	// Get user info from JWT middleware context
	userID, ok := c.Get("user_id").(int)
	if !ok {
		logger.Log.Warn().Msg("Failed to get user_id from context")
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Invalid user context"))
	}

	userRole, ok := c.Get("user_role").(string)
	if !ok {
		logger.Log.Warn().Int("user_id", userID).Msg("Failed to get user_role from context")
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Invalid user role"))
	}

	// Get document type and filename from URL params
	docType := c.Param("type")
	filename := c.Param("filename")

	// Validate document type
	validDocTypes := map[string]bool{
		"ktp":  true,
		"sim":  true,
		"stnk": true,
		"ktm":  true,
	}

	if !validDocTypes[strings.ToLower(docType)] {
		logger.Log.Warn().Str("doc_type", docType).Msg("Invalid document type requested")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_DOC_TYPE", "Invalid document type"))
	}

	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		logger.Log.Warn().
			Str("filename", filename).
			Int("user_id", userID).
			Msg("Directory traversal attempt detected")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_FILENAME", "Invalid filename"))
	}

	// Extract driver user ID from filename or path
	// Expected path structure: uploads/drivers/{userID}/{docType}/{filename}
	// We need to check all possible driver directories to find the file

	// For admin: allow access to any document
	// For driver: only allow access to own documents
	// For passenger: deny access

	if userRole == "PASSENGER" {
		logger.Log.Warn().
			Int("user_id", userID).
			Str("role", userRole).
			Msg("Passenger attempted to access driver document")
		return c.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", constants.ErrUnauthorizedAccess))
	}

	var filePath string

	if userRole == "ADMIN" {
		// Admin can access any document - search all driver directories
		driverDir := filepath.Join(h.uploadDir, "drivers")

		// Walk through all driver subdirectories to find the file
		found := false
		err := filepath.Walk(driverDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue walking even if there's an error
			}
			if !info.IsDir() && info.Name() == filename && strings.Contains(path, docType) {
				filePath = path
				found = true
				return filepath.SkipAll
			}
			return nil
		})

		if err != nil || !found {
			logger.Log.Warn().
				Str("filename", filename).
				Str("doc_type", docType).
				Msg("Document not found for admin")
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", constants.ErrDocumentNotFound))
		}
	} else if userRole == "DRIVER" {
		// Driver can only access their own documents
		driverIDStr := fmt.Sprintf("%d", userID)
		filePath = filepath.Join(h.uploadDir, "drivers", driverIDStr, docType, filename)
	} else {
		logger.Log.Warn().
			Int("user_id", userID).
			Str("role", userRole).
			Msg("Unknown role attempted to access document")
		return c.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", constants.ErrUnauthorizedAccess))
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Log.Warn().
			Str("file_path", filePath).
			Int("user_id", userID).
			Msg("Document file not found")
		return c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", constants.ErrDocumentNotFound))
	}

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".pdf":
		contentType = "application/pdf"
	}

	// Set headers and serve file
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")

	logger.Log.Info().
		Int("user_id", userID).
		Str("role", userRole).
		Str("doc_type", docType).
		Str("filename", filename).
		Msg("Document accessed successfully")

	return c.File(filePath)
}
