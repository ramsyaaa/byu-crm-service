package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/notification-one-signal/service"
	"byu-crm-service/modules/notification-one-signal/validation"

	"github.com/gofiber/fiber/v2"
)

type NotificationOneSignalHandler struct {
	notificationService service.NotificationOneSignalService
}

func NewNotificationOneSignalHandler(notificationService service.NotificationOneSignalService) *NotificationOneSignalHandler {
	return &NotificationOneSignalHandler{notificationService: notificationService}
}

func (h *NotificationOneSignalHandler) SubscribeNotification(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var successCode int = fiber.StatusOK

	req := new(validation.SubscribeNotificationRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateCreate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if req.Type != "Subscribe" && req.Type != "Unsubscribe" {
		errors := map[string]string{
			"type": "Type hanya boleh bernilai 'Subscribe' atau 'Unsubscribe'",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	uid := uint(userID)
	err := h.notificationService.CreateSubscribeNotification(&uid, req.SubscriptionID, req.Type)

	if err != nil {
		// Response
		response := helper.APIResponse("Internal Server Error", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Response
	response := helper.APIResponse("Success "+req.Type+" Notification", successCode, "success", nil)
	return c.Status(successCode).JSON(response)
}
