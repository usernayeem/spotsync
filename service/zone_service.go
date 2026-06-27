package service

import (
	"errors"

	"github.com/usernayeem/spotsync/dto"
	"github.com/usernayeem/spotsync/models"
	"github.com/usernayeem/spotsync/repository"
)

type ZoneService struct {
	zoneRepo *repository.ZoneRepository
}

func NewZoneService(zoneRepo *repository.ZoneRepository) *ZoneService {
	return &ZoneService{zoneRepo: zoneRepo}
}

func (s *ZoneService) CreateZone(req dto.CreateZoneRequest) (*models.ParkingZone, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.CreateZone(zone); err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *ZoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.GetAllZones()
	if err != nil {
		return nil, err
	}

	var responses []dto.ZoneResponse
	for _, zone := range zones {
		count, _ := s.zoneRepo.GetActiveReservationCount(zone.ID)
		
		responses = append(responses, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: zone.TotalCapacity - int(count),
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
			UpdatedAt:      zone.UpdatedAt,
		})
	}
	
	// Ensure we return an empty array instead of null in JSON if no zones exist
	if responses == nil {
		responses = []dto.ZoneResponse{}
	}
	
	return responses, nil
}

func (s *ZoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, errors.New("parking zone not found")
	}

	count, _ := s.zoneRepo.GetActiveReservationCount(zone.ID)

	response := &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(count),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}

	return response, nil
}

func (s *ZoneService) UpdateZone(id uint, req dto.UpdateZoneRequest) (*models.ParkingZone, error) {
	zone, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, errors.New("parking zone not found")
	}

	// Build only the fields that were actually provided
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.TotalCapacity != nil {
		updates["total_capacity"] = *req.TotalCapacity
	}
	if req.PricePerHour != nil {
		updates["price_per_hour"] = *req.PricePerHour
	}

	if err := s.zoneRepo.UpdateZone(zone, updates); err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *ZoneService) DeleteZone(id uint) error {
	_, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return errors.New("parking zone not found")
	}
	return s.zoneRepo.DeleteZone(id)
}

