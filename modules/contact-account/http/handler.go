package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/contact-account/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ContactAccountHandler struct {
	service service.ContactAccountService
}

func NewContactAccountHandler(service service.ContactAccountService) *ContactAccountHandler {
	return &ContactAccountHandler{service: service}
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
