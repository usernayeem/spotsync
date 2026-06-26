package models

import (
	"time"
)

// User represents the users table in the database
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // Omitted from JSON responses for security
	Role      string    `gorm:"type:varchar(20);default:'driver';not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
