package http

import (
	"byu-crm-service/modules/visit-checklist/service"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type VisitChecklistHandler struct {
	visitChecklistService service.VisitChecklistService
}

func NewVisitChecklistHandler(visitChecklistService service.VisitChecklistService) *VisitChecklistHandler {
	return &VisitChecklistHandler{visitChecklistService: visitChecklistService}
}

func (h *VisitChecklistHandler) GetAllVisitChecklist(c *fiber.Ctx) error {
	// Call service with filters
	visit_checklist, err := h.visitChecklistService.GetAllVisitChecklist()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch visit checklist",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"visit_checklist": visit_checklist,
	}

	response := helper.APIResponse("Get Visit Checklist Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
