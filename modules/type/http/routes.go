package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func TypeRoutes(router fiber.Router, handler *TypeHandler) {
	authRouter := router.Group("/types",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllTypes)
	authRouter.Get("/:id", handler.GetTypeByID)
	authRouter.Post("/", handler.CreateType)
	authRouter.Put("/:id", handler.UpdateType)
}
