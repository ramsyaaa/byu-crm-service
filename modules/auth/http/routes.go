package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router, handler *AuthHandler) {
	router.Post("/login", handler.Login)
	router.Post("/refresh", handler.Refresh)

	// Google OAuth routes
	router.Get("/google/login", handler.GoogleLogin)

	// Support both GET and POST for the callback
	router.Get("/callback/google", handler.GoogleCallback)
	router.Post("/callback/google", handler.GoogleCallback)

	authRouter := router.Group("/impersonate",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Post("/", handler.Impersonate)
}
