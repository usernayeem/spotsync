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
			Zone:         &res.Zone,
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

func (s *ReservationService) GetAllReservations() ([]models.Reservation, error) {
	return s.reservationRepo.GetAllReservations()
}
