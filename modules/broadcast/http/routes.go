package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func BroadcastRoutes(router fiber.Router, handler *BroadcastHandler) {
	authRouter := router.Group("/broadcasts",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/:id", handler.GetBroadcastByNotificationId)
	authRouter.Post("/", handler.CreateBroadcast)
}
