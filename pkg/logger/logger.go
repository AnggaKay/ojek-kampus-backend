package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// InitLogger initializes the global logger
func InitLogger(environment string) {
	// Configure based on environment
	if environment == "development" || environment == "" {
		// Pretty logging for development
		Log = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	} else {
		// JSON logging for production
		Log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Log.Info().Str("environment", environment).Msg("Logger initialized")
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	return &Log
}
