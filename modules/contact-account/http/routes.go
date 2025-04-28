package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ContactAccountRoutes(router fiber.Router, handler *ContactAccountHandler) {
	authRouter := router.Group("/contacts",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllContacts)
	authRouter.Get("/:id", handler.GetContactById)
	authRouter.Post("/", handler.CreateContact)
	authRouter.Put("/:id", handler.UpdateContact)
}
