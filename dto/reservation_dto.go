package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationZoneResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationResponse struct {
	ID           uint                     `json:"id"`
	UserID       uint                     `json:"user_id,omitempty"`
	LicensePlate string                   `json:"license_plate"`
	Status       string                   `json:"status"`
	Zone         *ReservationZoneResponse `json:"zone,omitempty"`
	CreatedAt    string                   `json:"created_at"`
}

type AdminUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AdminReservationResponse struct {
	ID           uint                     `json:"id"`
	UserID       uint                     `json:"user_id"`
	User         *AdminUserResponse       `json:"user,omitempty"`
	ZoneID       uint                     `json:"zone_id"`
	Zone         *ReservationZoneResponse `json:"zone,omitempty"`
	LicensePlate string                   `json:"license_plate"`
	Status       string                   `json:"status"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

