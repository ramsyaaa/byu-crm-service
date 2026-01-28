package http

import "github.com/gofiber/fiber/v2"

func PerformanceIndianaRoutes(app *fiber.App, handler *PerformanceIndianaHandler) {

	app.Post("/performance-indiana/import", handler.Import)
}
