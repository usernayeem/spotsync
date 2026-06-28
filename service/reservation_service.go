package service

import (
	"errors"

	"github.com/usernayeem/spotsync/dto"
	"github.com/usernayeem/spotsync/models"
	"github.com/usernayeem/spotsync/repository"
)

type ReservationService struct {
	reservationRepo *repository.ReservationRepository
}

func NewReservationService(reservationRepo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo}
}

func (s *ReservationService) ReserveSpot(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error) {
	// The repository handles the transaction and bottleneck logic
	return s.reservationRepo.CreateReservationTx(userID, req.ZoneID, req.LicensePlate)
}

func (s *ReservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.GetUserReservations(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ReservationResponse
	for _, res := range reservations {
		responses = append(responses, dto.ReservationResponse{
			ID:           res.ID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			Zone: &dto.ReservationZoneResponse{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			CreatedAt:    res.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	if responses == nil {
		responses = []dto.ReservationResponse{}
	}

	return responses, nil
}

func (s *ReservationService) CancelReservation(userID uint, reservationID uint) error {
	reservation, err := s.reservationRepo.GetReservationByID(reservationID)
	if err != nil {
		return errors.New("reservation not found")
	}

	// Verify the reservation belongs to the requester
	if reservation.UserID != userID {
		return errors.New("forbidden: you can only cancel your own reservations")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	return s.reservationRepo.UpdateReservationStatus(reservationID, "cancelled")
}

func (s *ReservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var responses []dto.AdminReservationResponse
	for _, res := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           res.ID,
			UserID:       res.UserID,
			User: &dto.AdminUserResponse{
				ID:    res.User.ID,
				Name:  res.User.Name,
				Email: res.User.Email,
				Role:  res.User.Role,
			},
			ZoneID: res.ZoneID,
			Zone: &dto.ReservationZoneResponse{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,
		})
	}

	if responses == nil {
		responses = []dto.AdminReservationResponse{}
	}

	return responses, nil
}
