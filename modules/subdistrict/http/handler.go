package http

import (
	"strconv"

	"byu-crm-service/modules/subdistrict/service"

	"github.com/gofiber/fiber/v2"
)

type SubdistrictHandler struct {
	service service.SubdistrictService
}

func NewSubdistrictHandler(service service.SubdistrictService) *SubdistrictHandler {
	return &SubdistrictHandler{service: service}
}

func (h *SubdistrictHandler) GetSubdistrictByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID parameter"})
	}

	Subdistrict, err := h.service.GetSubdistrictByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if Subdistrict == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Subdistrict not found"})
	}

	return c.JSON(Subdistrict)
}

func (h *SubdistrictHandler) GetSubdistrictByName(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Subdistrict name is required"})
	}

	Subdistrict, err := h.service.GetSubdistrictByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if Subdistrict == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Subdistrict not found"})
	}

	return c.JSON(Subdistrict)
}
