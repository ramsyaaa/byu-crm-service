package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/contact-account/service"
	"byu-crm-service/modules/contact-account/validation"
	socialMediaService "byu-crm-service/modules/social-media/service"
	"encoding/json"
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
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)

	// Call service with filters
	contacts, total, err := h.service.GetAllContacts(limit, paginate, page, filters, userRole, territoryID, userID)
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
	contact, err := h.service.CreateContact(reqMap)
	if err != nil {
		response := helper.APIResponse("Failed to create contact", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	_, _ = h.service.InsertContactAccountByContactID(reqMap, contact.ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Contact", contact.ID)

	// Return success response
	response := helper.APIResponse("Create Contact Succsesfully", fiber.StatusOK, "success", contact)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ContactAccountHandler) UpdateContact(c *fiber.Ctx) error {
	contactIDStr := c.Params("id")
	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid contact ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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

	contact, err := h.service.UpdateContact(reqMap, contactID)
	if err != nil {
		response := helper.APIResponse("Failed to update contact", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	_, _ = h.service.InsertContactAccountByContactID(reqMap, contact.ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Contact", contact.ID)

	// Return success response
	response := helper.APIResponse("Update Contact Succsesfully", fiber.StatusOK, "success", contact)
	return c.Status(fiber.StatusOK).JSON(response)
}
