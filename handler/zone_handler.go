package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/usernayeem/spotsync/dto"
	"github.com/usernayeem/spotsync/service"
)

type ZoneHandler struct {
	zoneService *service.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService *service.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validate:    validator.New(),
	}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
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

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create parking zone",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Parking zone created successfully",
		"data": map[string]interface{}{
			"id":             zone.ID,
			"name":           zone.Name,
			"type":           zone.Type,
			"total_capacity": zone.TotalCapacity,
			"price_per_hour": zone.PricePerHour,
			"created_at":     zone.CreatedAt,
			"updated_at":     zone.UpdatedAt,
		},
	})
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve parking zones",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    zones,
	})
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    zone,
	})
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	var req dto.UpdateZoneRequest
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

	zone, err := h.zoneService.UpdateZone(uint(id), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "parking zone not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone updated successfully",
		"data": map[string]interface{}{
			"id":             zone.ID,
			"name":           zone.Name,
			"type":           zone.Type,
			"total_capacity": zone.TotalCapacity,
			"price_per_hour": zone.PricePerHour,
			"created_at":     zone.CreatedAt,
			"updated_at":     zone.UpdatedAt,
		},
	})
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	if err := h.zoneService.DeleteZone(uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "parking zone not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone deleted successfully",
	})
}

