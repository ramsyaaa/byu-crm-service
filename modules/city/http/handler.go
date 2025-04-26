package http

import (
	"strconv"

	"byu-crm-service/helper"
	"byu-crm-service/modules/city/service"

	"github.com/gofiber/fiber/v2"
)

type CityHandler struct {
	service service.CityService
}

func NewCityHandler(service service.CityService) *CityHandler {
	return &CityHandler{service: service}
}

func (h *CityHandler) GetAllCities(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	cities, total, err := h.service.GetAllCities(limit, paginate, page, filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch cities",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"cities": cities,
		"total":  total,
		"page":   page,
	}

	response := helper.APIResponse("Get Cities Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CityHandler) GetCityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID parameter", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	city, err := h.service.GetCityByID(uint(id))
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	response := helper.APIResponse("Get City Successfully", fiber.StatusOK, "success", city)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CityHandler) GetCityByName(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "City name is required"})
	}

	city, err := h.service.GetCityByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if city == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "City not found"})
	}

	return c.JSON(city)
}
