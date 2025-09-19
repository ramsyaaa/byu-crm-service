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
	authRouter.Put("/profile", handler.UpdateUserProfile)
	authRouter.Get("/profile", handler.GetUserProfile)
	authRouter.Get("/resume-user/:id", handler.GetUserProfileResume)
	authRouter.Get("/resume-territory", handler.GetUsersResume)
	authRouter.Get("/:id", handler.GetUserByID)
	authRouter.Get("/", handler.GetAllUsers)
	authRouter.Post("/", handler.CreateUser)
	authRouter.Put("/:id", handler.UpdateUser)
}
