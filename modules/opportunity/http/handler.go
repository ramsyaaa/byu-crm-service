package http

import (
	"log"
	"byu-crm-service/models"
	"byu-crm-service/modules/opportunity/service"
	"byu-crm-service/modules/opportunity/validation"
	"context"
	"encoding/json"
	"fmt"
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
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf(fmt.Sprintf("Panic in Create Opportunity: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get user information from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateRequest)
	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateData(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
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

	// Create Opportunity with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var opportunity *models.Opportunity
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during opportunity creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		opportunity, serviceErr = h.opportunityService.CreateOpportunity(reqMap, userID)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to create opportunity: %v", serviceErr))
			response := helper.APIResponse("Failed to create opportunity: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Create Opportunity Succsesfully", fiber.StatusOK, "success", opportunity)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OpportunityHandler) UpdateOpportunity(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf(fmt.Sprintf("Panic in Update Opportunity: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get user information from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	opportunityIDStr := c.Params("id")
	if opportunityIDStr == "" {
		response := helper.APIResponse("Opportunity ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	opportunityID, err := strconv.Atoi(opportunityIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid Opportunity ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateRequest)
	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateData(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
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

	// Create Opportunity with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var opportunity *models.Opportunity
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during opportunity creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		opportunity, serviceErr = h.opportunityService.UpdateOpportunity(reqMap, userID, opportunityID)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to update opportunity: %v", serviceErr))
			response := helper.APIResponse("Failed to update opportunity: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Update Opportunity Succsesfully", fiber.StatusOK, "success", opportunity)
	return c.Status(fiber.StatusOK).JSON(response)
}
