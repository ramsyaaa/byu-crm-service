package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func VisitChecklistRoutes(router fiber.Router, handler *VisitChecklistHandler) {
	authRouter := router.Group("/visit-checklist",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllVisitChecklist)
}
