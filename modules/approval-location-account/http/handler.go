package http

import (
	"fmt"
	"strconv"

	"byu-crm-service/helper"
	accountService "byu-crm-service/modules/account/service"
	"byu-crm-service/modules/approval-location-account/service"

	"github.com/gofiber/fiber/v2"
)

type ApprovalLocationAccountHandler struct {
	service        service.ApprovalLocationAccountService
	accountService accountService.AccountService
}

func NewApprovalLocationAccountHandler(
	service service.ApprovalLocationAccountService, accountService accountService.AccountService) *ApprovalLocationAccountHandler {

	return &ApprovalLocationAccountHandler{
		service:        service,
		accountService: accountService}
}

func (h *ApprovalLocationAccountHandler) GetAllApprovalRequest(c *fiber.Ctx) error {
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
	approval_request, total, err := h.service.GetAllApprovalRequest(limit, paginate, page, filters, userRole, territoryID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch approval location accounts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"approval_request": approval_request,
		"total":            total,
		"page":             page,
	}

	response := helper.APIResponse("Get Approval Location Accounts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ApprovalLocationAccountHandler) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	account, err := h.service.FindByID(uint(id))
	if err != nil || account == nil {
		response := helper.APIResponse("Account not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Success get account", fiber.StatusOK, "success", account)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ApprovalLocationAccountHandler) HandleLocationApproval(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		errors := map[string]string{
			"id": "ID tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	status := c.Params("status")
	if status != "accept" && status != "reject" {
		errors := map[string]string{
			"status": "Status tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	RequestLocation, err := h.service.FindByID(uint(intID))
	if err != nil {
		response := helper.APIResponse("Failed to fetch request location", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if RequestLocation == nil {
		response := helper.APIResponse("Request not found", fiber.StatusOK, "error", nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}
	if status == "accept" {
		updateData := map[string]interface{}{
			"status": 1,
		}

		err = h.service.UpdateFields(uint(intID), updateData)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
		fmt.Println("RequestLocation.Longitude", &RequestLocation.Longitude)
		fmt.Println("RequestLocation.Latitude", &RequestLocation.Latitude)
		updatedLocationAccount := map[string]interface{}{
			"longitude": &RequestLocation.Longitude,
			"latitude":  &RequestLocation.Latitude,
		}

		err = h.accountService.UpdateFields(uint(*RequestLocation.AccountID), updatedLocationAccount)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
	} else if status == "reject" {
		updateData := map[string]interface{}{
			"status": 2,
		}

		err = h.service.UpdateFields(uint(intID), updateData)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
	}

	response := helper.APIResponse("Approval Location Account Successfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
