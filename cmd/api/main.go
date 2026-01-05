package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver PostgreSQL
)

func main() {
	// 1. Load file .env (Konfigurasi rahasia)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	// 2. Setup Database (PostgreSQL)
	// Format: postgres://user:password@host:port/dbname
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Tes koneksi DB
	if err := db.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}
	fmt.Println("✅ Database Connected Successfully!")

	// 3. Setup Framework Echo
	e := echo.New()

	// --- MIDDLEWARE PENTING UNTUK REACT ---
	e.Use(middleware.Logger())  // Supaya terlihat log request di terminal
	e.Use(middleware.Recover()) // Supaya server tidak crash kalau ada panic

	// CORS: Mengizinkan Frontend React akses API ini
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Di production ganti "*" dengan domain asli
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// 4. Routes (Endpoint API)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Ojek Kampus API 🚀",
			"status":  "active",
		})
	})

	// Contoh Endpoint Login (Dummy) untuk dites Frontend
	e.POST("/login", func(c echo.Context) error {
		// Simulasi respon login
		return c.JSON(http.StatusOK, map[string]interface{}{
			"token": "eyJhbGciOiJIUz...",
			"user": map[string]interface{}{
				"id":   1,
				"name": "Maba Teknik",
				"role": "PASSENGER",
			},
		})
	})

	// 5. Jalankan Server di Port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
