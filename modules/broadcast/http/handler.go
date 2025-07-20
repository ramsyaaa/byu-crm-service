package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/modules/broadcast/service"
	"byu-crm-service/modules/broadcast/validation"
	notificationService "byu-crm-service/modules/notification/service"
	roleService "byu-crm-service/modules/role/service"
	smsSender "byu-crm-service/modules/sms-sender/service"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
)

type BroadcastHandler struct {
	service             service.BroadcastService
	notificationService notificationService.NotificationService
	userService         userService.UserService
	roleService         roleService.RoleService
	smsSenderService    smsSender.SmsSenderService
}

func NewBroadcastHandler(
	service service.BroadcastService,
	notificationService notificationService.NotificationService,
	userService userService.UserService,
	roleService roleService.RoleService,
	smsSenderService smsSender.SmsSenderService) *BroadcastHandler {

	return &BroadcastHandler{
		service:             service,
		notificationService: notificationService,
		userService:         userService,
		roleService:         roleService,
		smsSenderService:    smsSenderService}
}

func (h *BroadcastHandler) GetBroadcastByNotificationId(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	// Get id from param
	idParam := c.Params("id")

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	notification, err := h.notificationService.GetByNotificationId(uint(id), uint(userID))
	if err != nil {
		response := helper.APIResponse("Broadcast not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Success get broadcast", fiber.StatusOK, "success", notification)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *BroadcastHandler) CreateBroadcast(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Create Broadcast: %v", r)
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Parse request body with error handling
	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if req.Type != "" && req.Type == "USER" {
		errors := validation.ValidateByUser(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.Type != "" && req.Type == "ROLE" {
		errors := validation.ValidateByRole(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		errors := map[string]string{
			"type": "Tipe tidak valid.",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)

	}

	// Create Account with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	requestBody := map[string]string{
		"title":        req.Title,
		"description":  req.Description,
		"message":      req.Description,
		"callback_url": "",
		"subject_type": "",
		"subject_id":   "",
	}

	if req.Type != "" && req.Type == "USER" {

		// Convert *[]string to []int
		var userIDs []int
		if req.UserID != nil {
			for _, idStr := range *req.UserID {
				idInt, convErr := strconv.Atoi(idStr)
				if convErr != nil {
					response := helper.APIResponse("Invalid user ID: "+idStr, fiber.StatusBadRequest, "error", nil)
					return c.Status(fiber.StatusBadRequest).JSON(response)
				}
				userIDs = append(userIDs, idInt)
			}
		}

		err := h.notificationService.AssignNotificationToUsers(requestBody, userIDs)
		if err != nil {
			response := helper.APIResponse("Error create notification", fiber.StatusBadRequest, "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = h.smsSenderService.AssignSmsToUsers(requestBody, userIDs)
		if err != nil {
			response := helper.APIResponse("Error sending SMS", fiber.StatusBadRequest, "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

	} else if req.Type != "" && req.Type == "ROLE" {
		roles, err := h.roleService.GetRoleByIDs(req.RoleID)
		if err != nil {
			response := helper.APIResponse("Error fetching roles", fiber.StatusInternalServerError, "error", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		roleNames := []string{}
		for _, role := range roles {
			roleNames = append(roleNames, role.Name)
		}

		err = h.notificationService.CreateNotification(requestBody, roleNames, userRole, territoryID, 0)
		if err != nil {
			response := helper.APIResponse("Error create notification", fiber.StatusBadRequest, "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = h.smsSenderService.CreateSms(requestBody, roleNames, userRole, territoryID, 0)
		if err != nil {
			response := helper.APIResponse("Error sending SMS", fiber.StatusBadRequest, "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Create broadcast Succsesfully", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
