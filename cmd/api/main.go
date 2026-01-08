package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AnggaKay/ojek-kampus-backend/internal/handler"
	"github.com/AnggaKay/ojek-kampus-backend/internal/middleware"
	"github.com/AnggaKay/ojek-kampus-backend/internal/repository"
	"github.com/AnggaKay/ojek-kampus-backend/internal/service"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	// Initialize database connection pool
	db, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, refreshTokenRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)

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
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// Protected routes
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.JWTAuth())
	authProtected.GET("/me", authHandler.GetProfile)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Server starting on port %s...\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET  /health")
	fmt.Println("   POST /api/auth/register/passenger")
	fmt.Println("   POST /api/auth/login")
	fmt.Println("   POST /api/auth/refresh")
	fmt.Println("   POST /api/auth/logout")
	fmt.Println("   GET  /api/auth/me (protected)")
	fmt.Println()

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
