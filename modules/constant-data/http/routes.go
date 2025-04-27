package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ConstantDataRoutes(router fiber.Router, handler *ConstantDataHandler) {
	authRouter := router.Group("/constant-data",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllConstants)
	authRouter.Post("/", handler.CreateConstant)
}
