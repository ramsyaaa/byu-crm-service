package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func AbsenceUserRoutes(router fiber.Router, handler *AbsenceUserHandler) {
	authRouter := router.Group("/absence-user",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllAbsenceUsers)
	authRouter.Get("/active-absence", handler.GetAbsenceActive)
	authRouter.Get("/:id", handler.GetAbsenceUserByID)
	authRouter.Post("/", handler.CreateAbsenceUser)
}
