package http

import "github.com/gofiber/fiber/v2"

func AuthRoutes(router fiber.Router, handler *AuthHandler) {
	router.Post("/login", handler.Login)

	// Google OAuth routes
	router.Get("/google/login", handler.GoogleLogin)

	// Support both GET and POST for the callback
	router.Get("/callback/google", handler.GoogleCallback)
	router.Post("/callback/google", handler.GoogleCallback)
}
