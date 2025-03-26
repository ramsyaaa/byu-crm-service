package http

import "github.com/gofiber/fiber/v2"

func AuthRoutes(app *fiber.App, handler *AuthHandler) {

	app.Post("/login", handler.Login)
}
