package http

import "github.com/gofiber/fiber/v2"

func PerformanceDigiposRoutes(app *fiber.App, handler *PerformanceDigiposHandler) {

	app.Post("/performance-digipos/import", handler.Import)
}
