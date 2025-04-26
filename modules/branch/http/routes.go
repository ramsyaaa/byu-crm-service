package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func BranchRoutes(router fiber.Router, handler *BranchHandler) {
	authRouter := router.Group("/branches",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllBranches)
	authRouter.Get("/:id", handler.GetBranchByID)
	authRouter.Post("/", handler.CreateBranch)
	authRouter.Put("/:id", handler.UpdateBranch)
}
