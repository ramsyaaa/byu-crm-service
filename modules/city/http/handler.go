package http

import (
	"byu-crm-service/modules/city/service"
	"byu-crm-service/modules/city/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type CityHandler struct {
	cityService service.CityService
}

func NewCityHandler(cityService service.CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

func (h *CityHandler) GetAllCities(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search": c.Query("search", ""),
	}

	// Parse integer and boolean values
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	cities, total, err := h.cityService.GetAllCities(filters, userRole, territoryID)
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
	}

	response := helper.APIResponse("Get Cities Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CityHandler) GetCityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid City ID",
			"error":   err.Error(),
		})
	}
	city, err := h.cityService.GetCityByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch city",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"city": city,
	}

	response := helper.APIResponse("Get City Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CityHandler) CreateCity(c *fiber.Ctx) error {
	req := new(validation.CreateCityRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateCreate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	city, err := h.cityService.CreateCity(&req.Name, req.ClusterID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("City created successful", fiber.StatusOK, "success", city)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CityHandler) UpdateCity(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid city ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateCityRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateUpdate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	currentCity, _ := h.cityService.GetCityByID(intID)
	if currentCity == nil {
		errors := map[string]string{
			"name": "City tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	city, err := h.cityService.UpdateCity(&req.Name, req.ClusterID, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("City updated successful", fiber.StatusOK, "success", city)
	return c.Status(fiber.StatusOK).JSON(response)
}
