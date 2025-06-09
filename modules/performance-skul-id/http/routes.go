package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func PerformanceSkulIdRoutes(router fiber.Router, handler *PerformanceSkulIdHandler) {

	authRouter := router.Group("/performance-skul-id",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)

	authRouter.Post("/import", handler.Import)
	authRouter.Post("/import-by-account/:id", handler.ImportByAccount)
	authRouter.Get("/", handler.GetAllSkulIds)
	authRouter.Post("/:account_id", handler.CreatePerformanceSkulID)
}
