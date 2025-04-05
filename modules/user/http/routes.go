package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, handler *UserHandler) {
	authRouter := router.Group("/users",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/profile", handler.GetUserProfile)
	authRouter.Get("/:id", handler.GetUserByID)
	authRouter.Get("/", handler.GetAllUsers)
}
