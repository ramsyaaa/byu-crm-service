package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func CommunicationRoutes(router fiber.Router, handler *CommunicationHandler) {
	authRouter := router.Group("/communications",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllCommunications)
	authRouter.Get("/:id", handler.GetCommunicationByID)
	authRouter.Post("/", handler.CreateCommunication)
	authRouter.Put("/:id", handler.UpdateCommunication)
}
