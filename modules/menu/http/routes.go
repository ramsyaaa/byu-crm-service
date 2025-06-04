package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func MenuRoutes(router fiber.Router, handler *MenuHandler) {
	authRouter := router.Group("/menus",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllMenus)
}
