package http

import (
	"strconv"

	"byu-crm-service/modules/city/service"

	"github.com/gofiber/fiber/v2"
)

type CityHandler struct {
	service service.CityService
}

func NewCityHandler(service service.CityService) *CityHandler {
	return &CityHandler{service: service}
}

func (h *CityHandler) GetCityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID parameter"})
	}

	city, err := h.service.GetCityByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if city == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "City not found"})
	}

	return c.JSON(city)
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
