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
	// authRouter.Get("/:id", handler.GetCategoryByID)
	// authRouter.Post("/", handler.CreateCategory)
	// authRouter.Put("/:id", handler.UpdateCategory)
}
