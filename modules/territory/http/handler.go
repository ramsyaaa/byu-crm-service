package http

import (
	"byu-crm-service/modules/territory/service"
	"fmt"
	"strconv"

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

func (h *TerritoryHandler) GetAllTerritoryResume(c *fiber.Ctx) error {
	var userRole string
	if queryUserRole := c.Query("user_role"); queryUserRole != "" {
		userRole = queryUserRole
	} else {
		local := c.Locals("user_role")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_role tidak ditemukan",
			})
		}
		userRoleStr, ok := local.(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_role is not a string",
			})
		}
		userRole = userRoleStr
	}
	var territoryID int
	if queryTerritoryID := c.Query("territory_id"); queryTerritoryID != "" {
		var err error
		territoryID, err = strconv.Atoi(queryTerritoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid territory_id parameter",
				"error":   err.Error(),
			})
		}
	} else {
		local := c.Locals("territory_id")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "territory_id tidak ditemukan",
			})
		}
		territoryIDStr, ok := local.(int)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "territory_id is not a int",
			})
		}
		territoryID = territoryIDStr
	}

	fmt.Println("userRole:", userRole)
	fmt.Println("territoryID:", territoryID)

	// Call service with filters
	// territories, total, err := h.service.GetAllTerritoryResume(userRole, territoryID)
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": "Failed to fetch territories resume",
	// 		"error":   err.Error(),
	// 	})
	// }

	// Return response
	// responseData := map[string]interface{}{
	// 	"territories": territories,
	// 	"total":       total,
	// }

	response := helper.APIResponse("Get Territories Resume Successfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
