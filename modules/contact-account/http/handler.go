package http

import (
	"log"
	"byu-crm-service/helper"
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/service"
	"byu-crm-service/modules/contact-account/validation"
	socialMediaService "byu-crm-service/modules/social-media/service"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ContactAccountHandler struct {
	service            service.ContactAccountService
	socialMediaService socialMediaService.SocialMediaService
}

func NewContactAccountHandler(service service.ContactAccountService, socialMediaService socialMediaService.SocialMediaService) *ContactAccountHandler {
	return &ContactAccountHandler{service: service, socialMediaService: socialMediaService}
}

func (h *ContactAccountHandler) GetAllContacts(c *fiber.Ctx) error {
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
	account_id, _ := strconv.Atoi(c.Query("account_id", "0"))
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	contacts, total, err := h.service.GetAllContacts(limit, paginate, page, filters, userRole, territoryID, account_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch contacts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"contacts": contacts,
		"total":    total,
		"page":     page,
	}

	response := helper.APIResponse("Get Contacts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ContactAccountHandler) GetContactById(c *fiber.Ctx) error {
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

	contact, err := h.service.FindByContactID(uint(id), userRole, uint(territoryID))
	if err != nil {
		response := helper.APIResponse("Contact not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	responseData := map[string]interface{}{
		"contact": contact,
	}

	response := helper.APIResponse("Success get contact", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ContactAccountHandler) CreateContact(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf(fmt.Sprintf("Panic in Create Contact: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

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

	// Request Validation
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

		if *req.Birthday != "" {
			_, err := time.Parse("2006-01-02", *req.Birthday)
			if err != nil {
				errors := map[string]string{
					"birthday": "Format tanggal lahir tidak benar",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}
	}

	// Create Account
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

	var contact *models.Contact
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during account creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		contact, serviceErr = h.service.CreateContact(reqMap)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to create account: %v", serviceErr))
			response := helper.APIResponse("Failed to create account: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	_, _ = h.service.InsertContactAccountByContactID(reqMap, contact.ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Contact", contact.ID)

	// Return success response
	response := helper.APIResponse("Create Contact Succsesfully", fiber.StatusOK, "success", contact)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ContactAccountHandler) UpdateContact(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf(fmt.Sprintf("Panic in Update Contact: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	contactIDStr := c.Params("id")
	if contactIDStr == "" {
		response := helper.APIResponse("Contact ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid contact ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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

	// Request Validation
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

		if *req.Birthday != "" {
			_, err := time.Parse("2006-01-02", *req.Birthday)
			if err != nil {
				errors := map[string]string{
					"birthday": "Format tanggal lahir tidak benar",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}
	}

	// Create Account
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

	var contact *models.Contact
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during contact update", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		contact, serviceErr = h.service.UpdateContact(reqMap, contactID)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to update contact: %v", serviceErr))
			response := helper.APIResponse("Failed to update contact: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	_, _ = h.service.InsertContactAccountByContactID(reqMap, contact.ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Contact", contact.ID)

	// Return success response
	response := helper.APIResponse("Update Contact Succsesfully", fiber.StatusOK, "success", contact)
	return c.Status(fiber.StatusOK).JSON(response)
}
