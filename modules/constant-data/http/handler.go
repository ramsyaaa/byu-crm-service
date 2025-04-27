package http

import (
	"byu-crm-service/modules/constant-data/service"
	"byu-crm-service/modules/constant-data/validation"
	"encoding/json"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type ConstantDataHandler struct {
	constantDataService service.ConstantDataService
}

func NewConstantDataHandler(ConstantDataService service.ConstantDataService) *ConstantDataHandler {
	return &ConstantDataHandler{constantDataService: ConstantDataService}
}

func (h *ConstantDataHandler) GetAllConstants(c *fiber.Ctx) error {
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
	type_constant := c.Query("type", "")
	other_group := c.Query("other_group", "")

	// Call service with filters
	constants, total, err := h.constantDataService.GetAllConstants(limit, paginate, page, filters, type_constant, other_group)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch constants",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"constants": constants,
		"total":     total,
		"page":      page,
	}

	response := helper.APIResponse("Get Constants Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ConstantDataHandler) CreateConstant(c *fiber.Ctx) error {
	req := new(validation.CreateConstantDataRequest)
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

	req.Value = strings.ToUpper(strings.TrimSpace(req.Value))
	req.Label = strings.ToUpper(strings.TrimSpace(req.Label))

	_, err := h.constantDataService.GetConstantByTypeAndValue(req.Type, req.Value)
	if err == nil {
		errors := map[string]string{
			"value": "Constant sudah tersedia",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	_ = json.Unmarshal(reqBytes, &reqMap)

	constant, err := h.constantDataService.CreateConstant(reqMap)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Constant created successful", fiber.StatusOK, "success", constant)
	return c.Status(fiber.StatusOK).JSON(response)
}
