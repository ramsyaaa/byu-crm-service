package http

import "github.com/gofiber/fiber/v2"

func CityRoutes(app *fiber.App, handler *CityHandler) {

	app.Get("/cities/:id", handler.GetCityByID)
	app.Get("/cities", handler.GetCityByName) // Query ?name=cityName
}
