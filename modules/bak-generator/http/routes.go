package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func BakGeneratorRoutes(router fiber.Router, handler *BakGeneratorHandler) {
	authRouter := router.Group("/bak",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Post("/create", handler.CreateBakGenerator)
	authRouter.Get("/:id", handler.GetBakByID)
	authRouter.Get("/", handler.GetAllBak)
}
