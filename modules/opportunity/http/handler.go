package http

import (
	"byu-crm-service/modules/opportunity/service"
	"strconv"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type OpportunityHandler struct {
	opportunityService service.OpportunityService
}

func NewOpportunityHandler(opportunityService service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{opportunityService: opportunityService}
}

func (h *OpportunityHandler) GetAllOpportunities(c *fiber.Ctx) error {
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
	opportunities, total, err := h.opportunityService.GetAllOpportunities(limit, paginate, page, filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch opportunities",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"opportunities": opportunities,
		"total":         total,
		"page":          page,
	}

	response := helper.APIResponse("Get Opportunities Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OpportunityHandler) GetOpportunityByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	// Get id from param
	idParam := c.Params("id")
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	opportunity, err := h.opportunityService.FindByOpportunityID(uint(id), userRole, uint(territoryID), uint(userID))
	if err != nil {
		response := helper.APIResponse("Opportunity not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	responseData := map[string]interface{}{
		"opportunity": opportunity,
	}

	response := helper.APIResponse("Success get Opportunity", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
