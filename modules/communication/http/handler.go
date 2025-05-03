package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/models"
	"byu-crm-service/modules/communication/service"
	"byu-crm-service/modules/communication/validation"
	opportunityService "byu-crm-service/modules/opportunity/service"

	"github.com/gofiber/fiber/v2"
)

type CommunicationHandler struct {
	service            service.CommunicationService
	opportunityService opportunityService.OpportunityService
}

func NewCommunicationHandler(service service.CommunicationService, opportunityService opportunityService.OpportunityService) *CommunicationHandler {
	return &CommunicationHandler{service: service, opportunityService: opportunityService}
}

func (h *CommunicationHandler) GetAllCommunications(c *fiber.Ctx) error {
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
	accountID, _ := strconv.Atoi(c.Query("account_id", "0"))

	// Call service with filters
	communications, total, err := h.service.GetAllCommunications(limit, paginate, page, filters, accountID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch communications",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"communications": communications,
		"total":          total,
		"page":           page,
	}

	response := helper.APIResponse("Get Communications Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CommunicationHandler) GetCommunicationByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	communication, err := h.service.FindByCommunicationID(uint(id))
	if err != nil {
		response := helper.APIResponse("Communication not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	responseData := map[string]interface{}{
		"communication": communication,
	}

	response := helper.APIResponse("Success get communication", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CommunicationHandler) CreateCommunication(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Create Communication: %v", r))
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
	req := new(validation.ValidateCreateRequest)
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
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if req.CommunicationType != "" && req.CommunicationType == "MENAWARKAN PROGRAM" {

		errors := validation.ValidateStatus(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		req.StatusCommunication = nil
	}

	if req.CheckOpportunity != nil && *req.CheckOpportunity == "1" && req.OpportunityName == nil {
		errors := map[string]string{
			"opportunity_name": "Nama opportunity wajib diisi",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create Account with context and error handling
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
			helper.LogError(c, fmt.Sprintf("Failed to marshal request: %v", marshalErr))
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
			helper.LogError(c, fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	reqMap["created_by"] = userID
	reqMap["opportunity_id"] = nil
	if reqMap["check_opportunity"] == "1" {
		reqMap["description"] = reqMap["note"]
		opportunity, err := h.opportunityService.CreateOpportunity(reqMap, userID)
		if err != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create opportunity: %v", err))
			response := helper.APIResponse("Failed to create opportunity: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		reqMap["opportunity_id"] = opportunity.ID
	}

	// Call service with timeout
	var communication *models.Communication
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during communication creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		communication, serviceErr = h.service.CreateCommunication(reqMap, userID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create communication: %v", serviceErr))
			response := helper.APIResponse("Failed to create communication: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Create Communication Succsesfully", fiber.StatusOK, "success", communication)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CommunicationHandler) UpdateCommunication(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Update Communication: %v", r))
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

	// Get and validate communication ID
	communicationIDStr := c.Params("id")
	if communicationIDStr == "" {
		response := helper.APIResponse("Communication ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	communicationID, err := strconv.Atoi(communicationIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid Communication ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateCreateRequest)
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
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if req.CommunicationType != "" && req.CommunicationType == "MENAWARKAN PROGRAM" {

		errors := validation.ValidateStatus(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		req.StatusCommunication = nil
	}

	if req.CheckOpportunity != nil && *req.CheckOpportunity == "1" && req.OpportunityName == nil {
		errors := map[string]string{
			"opportunity_name": "Nama opportunity wajib diisi",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create Account with context and error handling
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
			helper.LogError(c, fmt.Sprintf("Failed to marshal request: %v", marshalErr))
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
			helper.LogError(c, fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	reqMap["created_by"] = userID
	reqMap["opportunity_id"] = nil
	if reqMap["check_opportunity"] == "1" {
		reqMap["description"] = reqMap["note"]
		opportunity, err := h.opportunityService.CreateOpportunity(reqMap, userID)
		if err != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create opportunity: %v", err))
			response := helper.APIResponse("Failed to create opportunity: "+err.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		reqMap["opportunity_id"] = opportunity.ID
	}

	// Call service with timeout
	var communication *models.Communication
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during communication update", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		communication, serviceErr = h.service.UpdateCommunication(reqMap, userID, communicationID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to update communication: %v", serviceErr))
			response := helper.APIResponse("Failed to update communication: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Update Communication Succsesfully", fiber.StatusOK, "success", communication)
	return c.Status(fiber.StatusOK).JSON(response)
}
