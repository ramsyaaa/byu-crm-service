package http

import (
	"byu-crm-service/modules/region/service"
	"byu-crm-service/modules/region/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type RegionHandler struct {
	regionService service.RegionService
}

func NewRegionHandler(regionService service.RegionService) *RegionHandler {
	return &RegionHandler{regionService: regionService}
}

func (h *RegionHandler) GetAllRegions(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search": c.Query("search", ""),
	}

	// Parse integer and boolean values
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	regions, total, err := h.regionService.GetAllRegions(filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch regions",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"regions": regions,
		"total":   total,
	}

	response := helper.APIResponse("Get Regions Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegionHandler) GetRegionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid region ID",
			"error":   err.Error(),
		})
	}
	region, err := h.regionService.GetRegionByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch region",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"region": region,
	}

	response := helper.APIResponse("Get Region Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegionHandler) CreateRegion(c *fiber.Ctx) error {
	req := new(validation.CreateRegionRequest)
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

	existingRegion, _ := h.regionService.GetRegionByName(req.Name)
	if existingRegion != nil {
		errors := map[string]string{
			"name": "Nama regional sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	region, err := h.regionService.CreateRegion(&req.Name, req.AreaID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Region created successful", fiber.StatusOK, "success", region)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegionHandler) UpdateRegion(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid region ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateRegionRequest)
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

	currentRegion, _ := h.regionService.GetRegionByID(intID)
	if currentRegion == nil {
		errors := map[string]string{
			"name": "Regional tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	existingRegion, _ := h.regionService.GetRegionByName(req.Name)

	if existingRegion != nil && currentRegion.Name != req.Name {
		errors := map[string]string{
			"name": "Nama regional sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	region, err := h.regionService.UpdateRegion(&req.Name, req.AreaID, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Region updated successful", fiber.StatusOK, "success", region)
	return c.Status(fiber.StatusOK).JSON(response)
}
