package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usernayeem/spotsync/config"
	"github.com/usernayeem/spotsync/models"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize database
	db := config.InitDB()

	// Auto-Migrate Models
	err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check route
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
