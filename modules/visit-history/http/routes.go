package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func VisitHistoryRoutes(router fiber.Router, handler *VisitHistoryHandler) {
	authRouter := router.Group("/absence-user",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllAbsenceUsers)
	authRouter.Get("/:id", handler.GetAbsenceUserByID)
	authRouter.Post("/", handler.CreateAbsenceUser)
}
