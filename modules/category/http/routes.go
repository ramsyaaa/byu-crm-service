package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(router fiber.Router, handler *CategoryHandler) {
	authRouter := router.Group("/categories",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllCategories)
	authRouter.Get("/:id", handler.GetCategory)
	authRouter.Post("/", handler.CreateCategory)
	authRouter.Put("/:id", handler.UpdateCategory)
}
