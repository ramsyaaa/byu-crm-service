package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func YaeLeaderboardRoutes(router fiber.Router, handler *YaeLeaderboardHandler) {
	authRouter := router.Group("/yae-leaderboards",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllLeaderboards)
}
