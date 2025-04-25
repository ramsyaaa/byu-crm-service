package http

import (
	"byu-crm-service/modules/category/service"
	"strconv"

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

// func (h *CategoryHandler) GetFacultyByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	intID, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"message": "Invalid faculty ID",
// 			"error":   err.Error(),
// 		})
// 	}
// 	faculty, err := h.facultyService.GetFacultyByID(intID)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Failed to fetch faculty",
// 			"error":   err.Error(),
// 		})
// 	}

// 	// Return response
// 	responseData := map[string]interface{}{
// 		"faculty": faculty,
// 	}

// 	response := helper.APIResponse("Get Faculty Successfully", fiber.StatusOK, "success", responseData)
// 	return c.Status(fiber.StatusOK).JSON(response)
// }

// func (h *CategoryHandler) CreateFaculty(c *fiber.Ctx) error {
// 	req := new(validation.CreateFacultyRequest)
// 	if err := c.BodyParser(req); err != nil {
// 		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Request Validation
// 	errors := validation.ValidateCreate(req)
// 	if errors != nil {
// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	existingFaculty, _ := h.facultyService.GetFacultyByName(req.Name)
// 	if existingFaculty != nil {
// 		errors := map[string]string{
// 			"name": "Nama Fakultas sudah digunakan",
// 		}
// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	faculty, err := h.facultyService.CreateFaculty(&req.Name)
// 	if err != nil {
// 		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
// 		return c.Status(fiber.StatusUnauthorized).JSON(response)
// 	}

// 	// Response
// 	response := helper.APIResponse("Faculty created successful", fiber.StatusOK, "success", faculty)
// 	return c.Status(fiber.StatusOK).JSON(response)
// }

// func (h *CategoryHandler) UpdateFaculty(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	intID, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"message": "Invalid faculty ID",
// 			"error":   err.Error(),
// 		})
// 	}
// 	req := new(validation.UpdateFacultyRequest)
// 	if err := c.BodyParser(req); err != nil {
// 		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Request Validation
// 	errors := validation.ValidateUpdate(req)
// 	if errors != nil {
// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	currentFaculty, _ := h.facultyService.GetFacultyByID(intID)
// 	if currentFaculty == nil {
// 		errors := map[string]string{
// 			"name": "Fakultas tidak ditemukan",
// 		}
// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	existingFaculty, err := h.facultyService.GetFacultyByName(req.Name)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Failed to fetch faculty",
// 			"error":   err.Error(),
// 		})
// 	}

// 	if existingFaculty != nil && *currentFaculty.Name != req.Name {
// 		errors := map[string]string{
// 			"name": "Nama Fakultas sudah digunakan",
// 		}
// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	faculty, err := h.facultyService.UpdateFaculty(&req.Name, intID)
// 	if err != nil {
// 		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
// 		return c.Status(fiber.StatusUnauthorized).JSON(response)
// 	}

// 	// Response
// 	response := helper.APIResponse("Faculty updated successful", fiber.StatusOK, "success", faculty)
// 	return c.Status(fiber.StatusOK).JSON(response)
// }
