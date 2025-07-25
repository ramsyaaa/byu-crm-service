package http

import (
	"fmt"
	"strconv"

	"byu-crm-service/helper"
	accountService "byu-crm-service/modules/account/service"
	"byu-crm-service/modules/approval-location-account/service"
	notificationService "byu-crm-service/modules/notification/service"
	smsSenderService "byu-crm-service/modules/sms-sender/service"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
)

type ApprovalLocationAccountHandler struct {
	service             service.ApprovalLocationAccountService
	accountService      accountService.AccountService
	notificationService notificationService.NotificationService
	smsSenderService    smsSenderService.SmsSenderService
	userService         userService.UserService
}

func NewApprovalLocationAccountHandler(
	service service.ApprovalLocationAccountService, accountService accountService.AccountService,
	notificationService notificationService.NotificationService, smsSenderService smsSenderService.SmsSenderService,
	userService userService.UserService) *ApprovalLocationAccountHandler {

	return &ApprovalLocationAccountHandler{
		service:             service,
		accountService:      accountService,
		notificationService: notificationService,
		smsSenderService:    smsSenderService,
		userService:         userService,
	}
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

	filter_users := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	getUser, total, err := h.userService.GetAllUsers(0, false, 1, filter_users, []string{}, false, userRole, territoryID)

	if err != nil {
		response := helper.APIResponse("Failed to fetch users", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var userIDs []int
	for _, u := range getUser {
		userIDs = append(userIDs, int(u.ID))
	}

	// Call service with filters
	approval_request, total, err := h.service.GetAllApprovalRequest(limit, paginate, page, filters, userRole, territoryID, userIDs)
	if err != nil {
		response := helper.APIResponse("Failed to fetch approval location accounts", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
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
	userID := c.Locals("user_id").(int)

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	approvalRequest, err := h.service.FindByID(uint(id))
	if err != nil || approvalRequest == nil {
		response := helper.APIResponse("Approval request not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	_ = h.notificationService.MarkNotificationAsReadBySubjectID("App\\Models\\ApprovalLocationAccount", uint(approvalRequest.ID), userID)

	response := helper.APIResponse("Success get approval request", fiber.StatusOK, "success", approvalRequest)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ApprovalLocationAccountHandler) HandleLocationApproval(c *fiber.Ctx) error {
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)
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

		updatedLocationAccount := map[string]interface{}{
			"longitude": &RequestLocation.Longitude,
			"latitude":  &RequestLocation.Latitude,
		}

		err = h.accountService.UpdateFields(uint(*RequestLocation.AccountID), updatedLocationAccount)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		getAccount, err := h.accountService.FindByAccountID(uint(*RequestLocation.AccountID), userRole, uint(territoryID), uint(userID))
		if err != nil {
			response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		requestBody := map[string]string{
			"title":        "Pengajuan Perubahan Lokasi Diterima",
			"description":  fmt.Sprintf("Permintaan perubahan data lokasi account %s anda diterima.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/accounts?type=list"),
			"subject_type": "App\\Models\\ApprovalLocationAccount",
			"subject_id":   fmt.Sprintf("%d", RequestLocation.ID),
		}
		_ = h.notificationService.CreateNotification(requestBody, []string{}, userRole, territoryID, int(*RequestLocation.UserID))

		requestBody = map[string]string{
			"message":      fmt.Sprintf("Permintaan perubahan data lokasi account %s anda diterima.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/accounts?type=list"),
		}
		_ = h.smsSenderService.CreateSms(requestBody, []string{}, userRole, territoryID, int(*RequestLocation.UserID))

	} else if status == "reject" {
		updateData := map[string]interface{}{
			"status": 2,
		}

		err = h.service.UpdateFields(uint(intID), updateData)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		getAccount, err := h.accountService.FindByAccountID(uint(*RequestLocation.AccountID), userRole, uint(territoryID), uint(userID))
		if err != nil {
			response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		requestBody := map[string]string{
			"title":        "Pengajuan Perubahan Lokasi Ditolak",
			"description":  fmt.Sprintf("Permintaan perubahan data lokasi account %s anda ditolak.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/accounts?type=list"),
			"subject_type": "App\\Models\\ApprovalLocationAccount",
			"subject_id":   fmt.Sprintf("%d", RequestLocation.ID),
		}
		_ = h.notificationService.CreateNotification(requestBody, []string{}, userRole, territoryID, int(*RequestLocation.UserID))

		requestBody = map[string]string{
			"message":      fmt.Sprintf("Permintaan perubahan data lokasi account %s anda ditolak.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/accounts?type=list"),
		}
		_ = h.smsSenderService.CreateSms(requestBody, []string{}, userRole, territoryID, int(*RequestLocation.UserID))
	}

	response := helper.APIResponse("Approval Location Account Successfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
