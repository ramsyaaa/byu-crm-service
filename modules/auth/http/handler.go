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
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Response
	response := helper.APIResponse("Login successful", fiber.StatusOK, "success", fiber.Map{"token": token})
	return c.Status(fiber.StatusOK).JSON(response)
}

// GoogleLogin returns the Google OAuth URL
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.authService.GetGoogleOAuthURL()
	response := helper.APIResponse("Google OAuth URL", fiber.StatusOK, "success", fiber.Map{"url": url})
	return c.Status(fiber.StatusOK).JSON(response)
}

// GoogleCallback handles the callback from Google OAuth
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// Support both query parameter and JSON body for the code
	var code string

	// First check if code is in the query parameters
	code = c.Query("code")

	// If not in query, try to parse from request body
	if code == "" {
		req := new(validation.GoogleCallbackRequest)
		if err := c.BodyParser(req); err == nil {
			// Request Validation
			errors := validation.ValidateGoogleCallback(req)
			if errors != nil {
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
			code = req.Code
		}
	}

	// If still no code, return error
	if code == "" {
		response := helper.APIResponse("Invalid request: missing code", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Exchange the code for a token
	token, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Return the token
	response := helper.APIResponse("Google login successful", fiber.StatusOK, "success", fiber.Map{"token": token})
	return c.Status(fiber.StatusOK).JSON(response)
}
