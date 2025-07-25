package http

import (
	"byu-crm-service/modules/notification/service"
	"strconv"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetAllNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Call service with filters
	notifications, total, err := h.notificationService.GetAllNotifications(filters, limit, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch notifications",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"notifications": notifications,
		"total":         total,
	}

	response := helper.APIResponse("Get Notifications Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *NotificationHandler) GetNotificationById(c *fiber.Ctx) error {
	notificationID, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid notification ID",
			"error":   err.Error(),
		})
	}

	userID := c.Locals("user_id").(int)

	notification, err := h.notificationService.GetByNotificationId(uint(notificationID), uint(userID))
	if err != nil {
		response := helper.APIResponse("Notification not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	err = h.notificationService.MarkNotificationAsRead(uint(notificationID), uint(userID))
	if err != nil {
		response := helper.APIResponse("Failed to mark notification as read", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if notification == nil {
		response := helper.APIResponse("Notification not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Get Notification Successfully", fiber.StatusOK, "success", notification)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *NotificationHandler) MarkAllNotificationAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	err := h.notificationService.MarkAllNotificationsAsRead(userID)
	if err != nil {
		response := helper.APIResponse("Failed to mark all notifications as read", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := helper.APIResponse("All notifications marked as read successfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
