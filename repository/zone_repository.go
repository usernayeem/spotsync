package repository

import (
	"github.com/usernayeem/spotsync/models"
	"gorm.io/gorm"
)

type ZoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) CreateZone(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *ZoneRepository) GetAllZones() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *ZoneRepository) GetZoneByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetActiveReservationCount helps calculate available spots
func (r *ZoneRepository) GetActiveReservationCount(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}

// UpdateZone applies partial updates to an existing parking zone
func (r *ZoneRepository) UpdateZone(zone *models.ParkingZone, updates map[string]interface{}) error {
	return r.db.Model(zone).Updates(updates).Error
}

// DeleteZone permanently removes a parking zone by ID
func (r *ZoneRepository) DeleteZone(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

