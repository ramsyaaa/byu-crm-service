package http

import "github.com/gofiber/fiber/v2"

func PerformanceNamiRoutes(app *fiber.App, handler *PerformanceNamiHandler) {

	app.Post("/performance-nami/import", handler.Import)
}
