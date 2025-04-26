package http

import (
	"byu-crm-service/modules/area/service"
	"byu-crm-service/modules/area/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AreaHandler struct {
	areaService service.AreaService
}

func NewAreaHandler(areaService service.AreaService) *AreaHandler {
	return &AreaHandler{areaService: areaService}
}

func (h *AreaHandler) GetAllAreas(c *fiber.Ctx) error {
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
	areas, total, err := h.areaService.GetAllAreas(limit, paginate, page, filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch areas",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"areas": areas,
		"total": total,
		"page":  page,
	}

	response := helper.APIResponse("Get Areas Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AreaHandler) GetAreaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid area ID",
			"error":   err.Error(),
		})
	}
	area, err := h.areaService.GetAreaByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch area",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"area": area,
	}

	response := helper.APIResponse("Get Area Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AreaHandler) CreateArea(c *fiber.Ctx) error {
	req := new(validation.CreateAreaRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateCreate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	existingArea, _ := h.areaService.GetAreaByName(req.Name)
	if existingArea != nil {
		errors := map[string]string{
			"name": "Nama area sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	area, err := h.areaService.CreateArea(&req.Name)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Area created successful", fiber.StatusOK, "success", area)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AreaHandler) UpdateArea(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid area ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateAreaRequest)
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

	currentArea, _ := h.areaService.GetAreaByID(intID)
	if currentArea == nil {
		errors := map[string]string{
			"name": "Area tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	existingArea, _ := h.areaService.GetAreaByName(req.Name)

	if existingArea != nil && currentArea.Name != req.Name {
		errors := map[string]string{
			"name": "Nama Area sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	area, err := h.areaService.UpdateArea(&req.Name, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Area updated successful", fiber.StatusOK, "success", area)
	return c.Status(fiber.StatusOK).JSON(response)
}
