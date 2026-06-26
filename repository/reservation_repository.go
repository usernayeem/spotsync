package repository

import (
	"errors"

	"github.com/usernayeem/spotsync/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// CreateReservationTx handles the core concurrency logic using Transactions and Row-Level Locking
func (r *ReservationRepository) CreateReservationTx(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation *models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		// 1. Lock the parking zone row for update to prevent race conditions (The EV Bottleneck fix)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			return errors.New("parking zone not found")
		}

		// 2. Count active reservations for this specific zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check capacity bottleneck
		if activeCount >= int64(zone.TotalCapacity) {
			return errors.New("parking zone is at full capacity")
		}

		// 4. Create the reservation
		newReservation := models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&newReservation).Error; err != nil {
			return err
		}

		reservation = &newReservation
		return nil // Commits transaction
	})

	return reservation, err
}

func (r *ReservationRepository) GetUserReservations(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) GetReservationByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	return &reservation, err
}

func (r *ReservationRepository) UpdateReservationStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}
