package http

import (
	"byu-crm-service/modules/subdistrict/service"
	"byu-crm-service/modules/subdistrict/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type SubdistrictHandler struct {
	subdistrictService service.SubdistrictService
}

func NewSubdistrictHandler(subdistrictService service.SubdistrictService) *SubdistrictHandler {
	return &SubdistrictHandler{subdistrictService: subdistrictService}
}

func (h *SubdistrictHandler) GetAllSubdistricts(c *fiber.Ctx) error {
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
	subdistricts, total, err := h.subdistrictService.GetAllSubdistricts(limit, paginate, page, filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch subdistricts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"subdistricts": subdistricts,
		"total":        total,
		"page":         page,
	}

	response := helper.APIResponse("Get Subdistricts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *SubdistrictHandler) GetSubdistrictByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Subdistrict ID",
			"error":   err.Error(),
		})
	}
	subdistrict, err := h.subdistrictService.GetSubdistrictByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch Subdistrict",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"subdistrict": subdistrict,
	}

	response := helper.APIResponse("Get Subdistrict Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *SubdistrictHandler) CreateSubdistrict(c *fiber.Ctx) error {
	req := new(validation.CreateSubdistrictRequest)
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

	subdistrict, err := h.subdistrictService.CreateSubdistrict(&req.Name, req.CityID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Subdistrict created successful", fiber.StatusOK, "success", subdistrict)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *SubdistrictHandler) UpdateSubdistrict(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid subdistrict ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateSubdistrictRequest)
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

	currentSubdistrict, _ := h.subdistrictService.GetSubdistrictByID(intID)
	if currentSubdistrict == nil {
		errors := map[string]string{
			"name": "Subdistrict tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	subdistrict, err := h.subdistrictService.UpdateSubdistrict(&req.Name, req.CityID, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Subdistrict updated successful", fiber.StatusOK, "success", subdistrict)
	return c.Status(fiber.StatusOK).JSON(response)
}
