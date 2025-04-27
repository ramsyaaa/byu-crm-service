package http

import (
	"byu-crm-service/modules/category/service"
	"byu-crm-service/modules/category/validation"
	"encoding/json"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(CategoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: CategoryService}
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
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

	// Call service with filters
	categories, total, err := h.categoryService.GetAllCategories(limit, paginate, page, filters, module)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch categories",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"categories": categories,
		"total":      total,
		"page":       page,
	}

	response := helper.APIResponse("Get Categories Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		response := helper.APIResponse("ID not found", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	category, err := h.categoryService.GetCategoryByID(intID)
	if err != nil {
		response := helper.APIResponse("category not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Get Category Successfully", fiber.StatusOK, "success", category)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	req := new(validation.CreateCategoryRequest)
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

	_, err := h.categoryService.GetCategoryByNameAndModuleType(req.Name, *req.ModuleType)
	if err == nil {
		errors := map[string]string{
			"value": "Kategori sudah tersedia",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	_ = json.Unmarshal(reqBytes, &reqMap)

	category, err := h.categoryService.CreateCategory(reqMap)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Category created successful", fiber.StatusOK, "success", category)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
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

	req := new(validation.UpdateCategoryRequest)
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

	typeData, err := h.categoryService.GetCategoryByNameAndModuleType(req.Name, *req.ModuleType)
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

	category, err := h.categoryService.UpdateCategory(int(idUint), reqMap)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Category updated successful", fiber.StatusOK, "success", category)
	return c.Status(fiber.StatusOK).JSON(response)
}
