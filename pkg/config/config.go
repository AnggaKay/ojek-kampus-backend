package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	MaxConns int
	MinConns int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string
	Environment string
	Timezone    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", constants.DefaultMaxConns),
			MinConns: getEnvAsInt("DB_MIN_CONNS", constants.DefaultMinConns),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		Server: ServerConfig{
			Port:        getEnv("PORT", constants.DefaultPort),
			Environment: getEnv("ENVIRONMENT", "development"),
			Timezone:    getEnv("TZ", "Asia/Jakarta"),
		},
	}

	// Validate required fields
	if config.Database.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
