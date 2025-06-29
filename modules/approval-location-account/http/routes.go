package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ApprovalLocationAccountRoutes(router fiber.Router, handler *ApprovalLocationAccountHandler) {
	authRouter := router.Group("/approval-location-accounts",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllApprovalRequest)
	authRouter.Get("/:id", handler.GetById)
	authRouter.Post("/:id/:status", handler.HandleLocationApproval)
}
