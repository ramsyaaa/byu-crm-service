package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(router fiber.Router, handler *ProductHandler) {
	authRouter := router.Group("/products",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllProducts)
	authRouter.Get("/:id", handler.GetProductById)
	authRouter.Post("/", handler.CreateProduct)
	authRouter.Put("/:id", handler.UpdateProduct)
}
