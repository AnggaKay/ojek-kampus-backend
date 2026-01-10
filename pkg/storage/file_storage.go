package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/google/uuid"
)

// FileStorage defines interface for file storage operations
type FileStorage interface {
	Upload(file *multipart.FileHeader, userID int, docType string) (string, error)
	Delete(filePath string) error
	GetFullPath(relativePath string) string
	FileExists(relativePath string) bool
}

// LocalStorage implements FileStorage for local filesystem
type LocalStorage struct {
	baseDir string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(baseDir string) FileStorage {
	return &LocalStorage{
		baseDir: baseDir,
	}
}

// Upload uploads a file to local storage
func (s *LocalStorage) Upload(file *multipart.FileHeader, userID int, docType string) (string, error) {
	// Validate file
	if err := validateFile(file); err != nil {
		return "", err
	}

	// Generate secure filename
	filename := generateSecureFilename(file.Filename)

	// Create directory structure: uploads/drivers/{userID}/{docType}/
	relativeDir := filepath.Join("drivers", fmt.Sprintf("%d", userID), strings.ToLower(docType))
	fullDir := filepath.Join(s.baseDir, relativeDir)

	// Create directory if not exists
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		logger.Log.Error().Err(err).Str("dir", fullDir).Msg("Failed to create directory")
		return "", fmt.Errorf(constants.ErrFailedToUploadFile + ": cannot create directory")
	}

	// Full file path
	relativePath := filepath.Join(relativeDir, filename)
	fullPath := filepath.Join(s.baseDir, relativePath)

	// Open source file
	src, err := file.Open()
	if err != nil {
		logger.Log.Error().Err(err).Str("file", file.Filename).Msg("Failed to open uploaded file")
		return "", fmt.Errorf(constants.ErrFailedToUploadFile + ": cannot open file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Log.Error().Err(err).Str("path", fullPath).Msg("Failed to create file")
		return "", fmt.Errorf(constants.ErrFailedToUploadFile + ": cannot create file")
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		logger.Log.Error().Err(err).Str("path", fullPath).Msg("Failed to write file")
		return "", fmt.Errorf(constants.ErrFailedToUploadFile + ": cannot write file")
	}

	logger.Log.Info().Str("path", relativePath).Int("user_id", userID).Str("doc_type", docType).Msg("File uploaded successfully")

	// Return relative path (for database storage)
	return relativePath, nil
}

// Delete deletes a file from local storage
func (s *LocalStorage) Delete(relativePath string) error {
	fullPath := filepath.Join(s.baseDir, relativePath)

	if err := os.Remove(fullPath); err != nil {
		logger.Log.Error().Err(err).Str("path", fullPath).Msg("Failed to delete file")
		return err
	}

	logger.Log.Info().Str("path", relativePath).Msg("File deleted successfully")
	return nil
}

// GetFullPath returns full filesystem path from relative path
func (s *LocalStorage) GetFullPath(relativePath string) string {
	return filepath.Join(s.baseDir, relativePath)
}

// FileExists checks if file exists
func (s *LocalStorage) FileExists(relativePath string) bool {
	fullPath := s.GetFullPath(relativePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// validateFile validates uploaded file
func validateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > constants.MaxFileSize {
		return fmt.Errorf(constants.ErrFileTooLarge+" (max %d MB)", constants.MaxFileSize/(1024*1024))
	}

	// Open file to check MIME type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("cannot open file")
	}
	defer src.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("cannot read file")
	}

	// Detect MIME type using magic bytes
	mimeType := http.DetectContentType(buffer)

	// Check if MIME type is allowed
	allowedTypes := []string{"image/jpeg", "image/png", "application/pdf"}
	isAllowed := false
	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf(constants.ErrInvalidFileType+": %s (allowed: JPG, PNG, PDF)", mimeType)
	}

	return nil
}

// generateSecureFilename generates a secure unique filename
func generateSecureFilename(originalFilename string) string {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(originalFilename))

	// Generate UUID
	id := uuid.New().String()

	// Generate timestamp
	timestamp := time.Now().Unix()

	// Format: {timestamp}_{uuid}.{ext}
	return fmt.Sprintf("%d_%s%s", timestamp, id, ext)
}

// SanitizeFilename removes dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators and special characters
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	return filename
}
