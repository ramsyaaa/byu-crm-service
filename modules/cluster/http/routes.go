package http

import "github.com/gofiber/fiber/v2"

func ClusterRoutes(app *fiber.App, handler *ClusterHandler) {

	app.Get("/clusters/:id", handler.GetClusterByID)
	app.Get("/clusters", handler.GetClusterByName) // Query ?name=ClusterName
}
