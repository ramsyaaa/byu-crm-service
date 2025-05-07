package http

import (
	"byu-crm-service/modules/territory/service"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type TerritoryHandler struct {
	service service.TerritoryService
}

func NewTerritoryHandler(service service.TerritoryService) *TerritoryHandler {
	return &TerritoryHandler{service: service}
}

func (h *TerritoryHandler) GetAllTerritories(c *fiber.Ctx) error {

	// Call service with filters
	territories, total, err := h.service.GetAllTerritories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch territories",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"territories": territories,
		"total":       total,
	}

	response := helper.APIResponse("Get Territories Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
