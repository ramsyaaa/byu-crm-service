package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func AreaRoutes(router fiber.Router, handler *AreaHandler) {
	authRouter := router.Group("/areas",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllAreas)
	authRouter.Get("/:id", handler.GetAreaByID)
	authRouter.Post("/", handler.CreateArea)
	authRouter.Put("/:id", handler.UpdateArea)
}
