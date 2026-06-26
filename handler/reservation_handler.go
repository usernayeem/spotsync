package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/usernayeem/spotsync/dto"
	"github.com/usernayeem/spotsync/service"
)

type ReservationHandler struct {
	reservationService *service.ReservationService
	validate           *validator.Validate
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validate:           validator.New(),
	}
}

func (h *ReservationHandler) ReserveSpot(c echo.Context) error {
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	userID := uint(c.Get("userID").(float64)) // JWT claims parses numbers as float64

	reservation, err := h.reservationService.ReserveSpot(userID, req)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data": map[string]interface{}{
			"id":            reservation.ID,
			"zone_id":       reservation.ZoneID,
			"license_plate": reservation.LicensePlate,
			"status":        reservation.Status,
			"created_at":    reservation.CreatedAt,
		},
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := uint(c.Get("userID").(float64))

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve reservations",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    reservations,
	})
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid reservation ID",
		})
	}

	userID := uint(c.Get("userID").(float64))

	err = h.reservationService.CancelReservation(userID, uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "forbidden: you can only cancel your own reservations" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "reservation not found" {
			statusCode = http.StatusNotFound
		}
		
		return c.JSON(statusCode, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve reservations",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "All reservations retrieved successfully",
		"data":    reservations,
	})
}
