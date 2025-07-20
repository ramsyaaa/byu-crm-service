package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func NotificationRoutes(router fiber.Router, handler *NotificationHandler) {
	authRouter := router.Group("/notifications",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllNotifications)
	authRouter.Post("/mark-all-notification-as-read", handler.MarkAllNotificationAsRead)
	authRouter.Get("/:id", handler.GetNotificationById)
}
