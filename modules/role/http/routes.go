package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(router fiber.Router, handler *RoleHandler) {
	authRouter := router.Group("/roles",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllRoles)
	authRouter.Get("/:id", handler.GetRoleByID)
	authRouter.Post("/", handler.CreateRole)
	authRouter.Put("/:id", handler.UpdateRole)
}
