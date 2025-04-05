package http

import (
	"byu-crm-service/modules/auth/service"
	"byu-crm-service/modules/auth/validation"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login Handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(validation.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateLogin(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get Token
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Login successful", fiber.StatusOK, "success", fiber.Map{"token": token})
	return c.Status(fiber.StatusOK).JSON(response)
}
