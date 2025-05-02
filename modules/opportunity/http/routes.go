package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func OpportunityRoutes(router fiber.Router, handler *OpportunityHandler) {
	authRouter := router.Group("/opportunities",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllOpportunities)
	authRouter.Get("/:id", handler.GetOpportunityByID)
	authRouter.Post("/", handler.CreateOpportunity)
	authRouter.Put("/:id", handler.UpdateOpportunity)
}
