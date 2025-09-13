package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func NotificationOneSignalRoutes(router fiber.Router, handler *NotificationOneSignalHandler) {
	authRouter := router.Group("/notification-subscription",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Post("/", handler.SubscribeNotification)
}
