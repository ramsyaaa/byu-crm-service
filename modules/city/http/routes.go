package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func CityRoutes(router fiber.Router, handler *CityHandler) {
	authRouter := router.Group("/cities",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllCities)
	authRouter.Get("/:id", handler.GetCityByID)
	authRouter.Post("/", handler.CreateCity)
	authRouter.Put("/:id", handler.UpdateCity)
}
