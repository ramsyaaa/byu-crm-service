package http

import "github.com/gofiber/fiber/v2"

func PerformanceSkulIdRoutes(app *fiber.App, handler *PerformanceSkulIdHandler) {

	app.Post("/performance-skul-id/import", handler.Import)
}
