package main

import (
	"fmt"
	"log"

	"github.com/AnggaKay/ojek-kampus-backend/internal/handler"
	"github.com/AnggaKay/ojek-kampus-backend/internal/middleware"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/internal/service"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/config"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/database"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/storage"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Server.Environment)
	logger.Log.Info().Str("environment", cfg.Server.Environment).Msg("Starting Ojek Kampus Backend")

	// Initialize database connection pool
	db, err := database.NewPostgresPool()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	logger.Log.Info().Msg("Database connection established")

	// Initialize file storage
	fileStorage := storage.NewLocalStorage(constants.UploadDirectory)
	logger.Log.Info().Str("upload_dir", constants.UploadDirectory).Msg("File storage initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	passengerRepo := repository.NewPassengerRepository(db)
	driverRepo := repository.NewDriverRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, passengerRepo, driverRepo, refreshTokenRepo)
	driverService := service.NewDriverService(userRepo, driverRepo, refreshTokenRepo, fileStorage)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	driverHandler := handler.NewDriverHandler(driverService)
	documentHandler := handler.NewDocumentHandler(constants.UploadDirectory)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Set custom validator
	e.Validator = middleware.NewValidator()

	// Global middlewares
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	// Health check
	e.GET("/health", healthHandler.Check)
	e.GET("/", healthHandler.Check)

	// API v1 routes
	api := e.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register/passenger", authHandler.RegisterPassenger)
	auth.POST("/register/driver", driverHandler.RegisterDriver)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// Protected routes
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.JWTAuth())
	authProtected.GET("/me", authHandler.GetProfile)

	// Document routes (protected - requires authentication)
	documents := api.Group("/documents")
	documents.Use(middleware.JWTAuth())
	documents.GET("/:type/:filename", documentHandler.GetDocument)

	// Start server
	logger.Log.Info().Str("port", cfg.Server.Port).Msg("Server starting")
	fmt.Printf("\n🚀 Server starting on port %s...\n", cfg.Server.Port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET  /health")
	fmt.Println("   POST /api/auth/register/passenger")
	fmt.Println("   POST /api/auth/register/driver (multipart/form-data)")
	fmt.Println("   POST /api/auth/login")
	fmt.Println("   POST /api/auth/refresh")
	fmt.Println("   POST /api/auth/logout")
	fmt.Println("   GET  /api/auth/me (protected)")
	fmt.Println("   GET  /api/documents/:type/:filename (protected)")
	fmt.Println()

	if err := e.Start(":" + cfg.Server.Port); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
