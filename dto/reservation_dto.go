package dto

import "github.com/usernayeem/spotsync/models"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationResponse struct {
	ID           uint                `json:"id"`
	UserID       uint                `json:"user_id,omitempty"`
	LicensePlate string              `json:"license_plate"`
	Status       string              `json:"status"`
	Zone         *models.ParkingZone `json:"zone,omitempty"`
	CreatedAt    string              `json:"created_at"`
}
