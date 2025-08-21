package http

import (
	"byu-crm-service/modules/template-file/service"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type TemplateFileHandler struct {
	templateFileService service.TemplateFileService
}

func NewTemplateFileHandler(TemplateFileService service.TemplateFileService) *TemplateFileHandler {
	return &TemplateFileHandler{templateFileService: TemplateFileService}
}

func (h *TemplateFileHandler) GetAllFiles(c *fiber.Ctx) error {

	type_file := c.Query("type", "")

	if type_file == "" {
		return c.Status(fiber.StatusBadRequest).JSON(helper.APIResponse("Type file is required", fiber.StatusBadRequest, "error", nil))
	}

	// Call service with filters
	files := h.templateFileService.GetAllTemplateFiles(type_file)

	// Return response
	responseData := map[string]interface{}{
		"files": files,
	}

	response := helper.APIResponse("Get Files Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
