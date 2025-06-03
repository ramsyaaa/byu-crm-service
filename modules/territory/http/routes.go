package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func TerritoryRoutes(router fiber.Router, handler *TerritoryHandler) {
	authRouter := router.Group("/territories",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllTerritories)
	authRouter.Get("/resume", handler.GetAllTerritoryResume)
}
