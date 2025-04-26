package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegionRoutes(router fiber.Router, handler *RegionHandler) {
	authRouter := router.Group("/regions",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllRegions)
	authRouter.Get("/:id", handler.GetRegionByID)
	authRouter.Post("/", handler.CreateRegion)
	authRouter.Put("/:id", handler.UpdateRegion)
}
