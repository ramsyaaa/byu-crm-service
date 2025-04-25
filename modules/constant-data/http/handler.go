package http

import (
	"byu-crm-service/modules/constant-data/service"
	"strconv"

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

	// Call service with filters
	constants, total, err := h.constantDataService.GetAllConstants(limit, paginate, page, filters, type_constant)
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
