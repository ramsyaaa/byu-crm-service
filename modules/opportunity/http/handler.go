package http

import (
	"byu-crm-service/modules/opportunity/service"
	"byu-crm-service/modules/opportunity/validation"
	"encoding/json"
	"strconv"
	"time"

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

func (h *OpportunityHandler) CreateOpportunity(c *fiber.Ctx) error {
	req := new(validation.ValidateRequest)
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

	if *req.OpenDate != "" {
		_, err := time.Parse("2006-01-02", *req.OpenDate)
		if err != nil {
			errors := map[string]string{
				"open_date": "Format tanggal open tidak benar",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if *req.CloseDate != "" {
		_, err := time.Parse("2006-01-02", *req.CloseDate)
		if err != nil {
			errors := map[string]string{
				"close_date": "Format tanggal close tidak benar",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	// Create Account
	reqMap := make(map[string]interface{})

	// Melakukan marshal dan menangani error
	reqBytes, err := json.Marshal(req)
	if err != nil {
		response := helper.APIResponse("Failed to marshal request", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Melakukan unmarshal
	err = json.Unmarshal(reqBytes, &reqMap)
	if err != nil {
		response := helper.APIResponse("Failed to unmarshal request", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// userID := c.Locals("user_id").(int)

	// opportunity, err := h.opportunityService.CreateOpportunity(reqMap, userID)
	// if err != nil {
	// 	response := helper.APIResponse("Failed to create opportunity", fiber.StatusInternalServerError, "error", err.Error())
	// 	return c.Status(fiber.StatusInternalServerError).JSON(response)
	// }

	// Return success response
	response := helper.APIResponse("Create Opportunity Succsesfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
