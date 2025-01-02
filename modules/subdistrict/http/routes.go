package http

import "github.com/gofiber/fiber/v2"

func SubdistrictRoutes(app *fiber.App, handler *SubdistrictHandler) {

	app.Get("/subdistricts/:id", handler.GetSubdistrictByID)
	app.Get("/subdistricts", handler.GetSubdistrictByName) // Query ?name=SubdistrictName
}
