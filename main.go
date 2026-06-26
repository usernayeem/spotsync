package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usernayeem/spotsync/config"
	"github.com/usernayeem/spotsync/handler"
	"github.com/usernayeem/spotsync/middlewares"
	"github.com/usernayeem/spotsync/models"
	"github.com/usernayeem/spotsync/repository"
	"github.com/usernayeem/spotsync/service"
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

	// Dependency Injection for Auth Module
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Dependency Injection for Zone Module
	zoneRepo := repository.NewZoneRepository(db)
	zoneService := service.NewZoneService(zoneRepo)
	zoneHandler := handler.NewZoneHandler(zoneService)

	// Dependency Injection for Reservation Module
	reservationRepo := repository.NewReservationRepository(db)
	reservationService := service.NewReservationService(reservationRepo)
	reservationHandler := handler.NewReservationHandler(reservationService)

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

	// API Routes Group
	api := e.Group("/api/v1")
	
	// Auth Routes
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Zone Routes
	zoneGroup := api.Group("/zones")
	zoneGroup.GET("", zoneHandler.GetAllZones)       // Public
	zoneGroup.GET("/:id", zoneHandler.GetZoneByID)   // Public
	
	// Create zone requires Auth and Admin role
	zoneGroup.POST("", zoneHandler.CreateZone, middlewares.RequireAuth, middlewares.RequireAdmin)

	// Reservation Routes
	reservationGroup := api.Group("/reservations")
	// Must be logged in to access these
	reservationGroup.Use(middlewares.RequireAuth)
	
	reservationGroup.POST("", reservationHandler.ReserveSpot)
	reservationGroup.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservationGroup.DELETE("/:id", reservationHandler.CancelReservation)
	
	// Admin only
	reservationGroup.GET("", reservationHandler.GetAllReservations, middlewares.RequireAdmin)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
