package http

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/type/service"
	"byu-crm-service/modules/type/validation"
	"encoding/json"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type TypeHandler struct {
	typeService service.TypeService
}

func NewTypeHandler(TypeService service.TypeService) *TypeHandler {
	return &TypeHandler{typeService: TypeService}
}

func (h *TypeHandler) GetAllTypes(c *fiber.Ctx) error {
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
	module := c.Query("module_type", "")
	category_name := c.Query("category_name", "")
	categoryNames := []string{}
	if category_name != "" {
		categoryNames = strings.Split(category_name, ",")
	}

	// Call service with filters
	types, total, err := h.typeService.GetAllTypes(limit, paginate, page, filters, module, categoryNames)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch types",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"types": types,
		"total": total,
		"page":  page,
	}

	response := helper.APIResponse("Get Types Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *TypeHandler) GetType(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		response := helper.APIResponse("ID not found", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var (
		typeData models.Type // ganti dengan model Type kamu
		err      error
	)

	intID, convErr := strconv.Atoi(id)
	if convErr == nil {
		// Jika bisa di-convert ke integer, cari berdasarkan ID
		typeData, err = h.typeService.GetTypeByID(intID)
	} else {
		// Kalau tidak bisa di-convert ke integer, cari berdasarkan Name
		typeData, err = h.typeService.GetTypeByName(id)
	}

	if err != nil {
		response := helper.APIResponse("Type not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	responseData := map[string]interface{}{
		"type": typeData,
	}

	response := helper.APIResponse("Get Type Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *TypeHandler) CreateType(c *fiber.Ctx) error {
	req := new(validation.CreateTypeRequest)
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

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	categoryID, err := strconv.Atoi(req.CategoryID)
	if err != nil {
		categoryID = 0
	}
	_, err = h.typeService.GetTypeByNameAndModuleType(req.Name, *req.ModuleType, categoryID)
	if err == nil {
		errors := map[string]string{
			"value": "Tipe sudah tersedia",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	_ = json.Unmarshal(reqBytes, &reqMap)

	typeData, err := h.typeService.CreateType(reqMap)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Type created successful", fiber.StatusOK, "success", typeData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *TypeHandler) UpdateType(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		response := helper.APIResponse("ID not found", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req := new(validation.UpdateTypeRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateUpdate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	categoryID, err := strconv.Atoi(req.CategoryID)
	if err != nil {
		categoryID = 0
	}
	typeData, err := h.typeService.GetTypeByNameAndModuleType(req.Name, *req.ModuleType, categoryID)
	if err == nil {
		if typeData.ID != uint(idUint) {
			errors := map[string]string{
				"value": "Tipe sudah tersedia",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	_ = json.Unmarshal(reqBytes, &reqMap)

	typeData, err = h.typeService.UpdateType(int(idUint), reqMap)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Type updated successful", fiber.StatusOK, "success", typeData)
	return c.Status(fiber.StatusOK).JSON(response)
}
