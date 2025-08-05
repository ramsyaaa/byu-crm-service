package http

import (
	"byu-crm-service/modules/branch/service"
	"byu-crm-service/modules/branch/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type BranchHandler struct {
	branchService service.BranchService
}

func NewBranchHandler(branchService service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

func (h *BranchHandler) GetAllBranches(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search": c.Query("search", ""),
	}

	withGeo := strings.ToLower(c.Query("with_geo", "0")) == "1"

	// Parse integer and boolean values
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	branches, total, err := h.branchService.GetAllBranches(filters, userRole, territoryID, withGeo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch branches",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"branches": branches,
		"total":    total,
	}

	response := helper.APIResponse("Get Branches Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BranchHandler) GetBranchByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid branch ID",
			"error":   err.Error(),
		})
	}
	branch, err := h.branchService.GetBranchByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch branch",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"branch": branch,
	}

	response := helper.APIResponse("Get Branch Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BranchHandler) CreateBranch(c *fiber.Ctx) error {
	req := new(validation.CreateBranchRequest)
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

	branch, err := h.branchService.CreateBranch(&req.Name, req.RegionID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Branch created successful", fiber.StatusOK, "success", branch)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BranchHandler) UpdateBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid branch ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateBranchRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateUpdate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	currentBranch, _ := h.branchService.GetBranchByID(intID)
	if currentBranch == nil {
		errors := map[string]string{
			"name": "Branch tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	branch, err := h.branchService.UpdateBranch(&req.Name, req.RegionID, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Branch updated successful", fiber.StatusOK, "success", branch)
	return c.Status(fiber.StatusOK).JSON(response)
}
