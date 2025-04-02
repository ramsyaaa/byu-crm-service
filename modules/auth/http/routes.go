package http

import "github.com/gofiber/fiber/v2"

func AuthRoutes(router fiber.Router, handler *AuthHandler) {
	router.Post("/login", handler.Login)
}
