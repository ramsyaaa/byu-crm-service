package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func PermissionRoutes(router fiber.Router, handler *PermissionHandler) {
	authRouter := router.Group("/permissions",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllPermissions)
	authRouter.Get("/:id", handler.GetPermissionByID)
	authRouter.Post("/", handler.CreatePermission)
	authRouter.Put("/:id", handler.UpdatePermission)
}
