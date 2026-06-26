package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection using GORM
func InitDB() *gorm.DB {
	// If DATABASE_URL is provided (like in NeonDB), use it directly
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)
	} else {
		// NeonDB connection strings might need an extra query parameter for older pgbouncer depending on setup,
		// but the default provided by NeonDB usually works flawlessly with GORM.
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established successfully")
	return db
}
