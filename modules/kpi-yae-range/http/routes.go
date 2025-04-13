package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func KpiYaeRangeRoutes(router fiber.Router, handler *KpiYaeRangeHandler) {
	authRouter := router.Group("/kpi-yae",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/current", handler.GetCurrentKpiYaeRanges)
	authRouter.Get("/performance", handler.GetPerformanceUser)
	authRouter.Post("/", handler.CreateKpiYaeRange)
}
