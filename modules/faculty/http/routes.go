package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func FacultyRoutes(router fiber.Router, handler *FacultyHandler) {
	authRouter := router.Group("/faculties",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllFaculties)
	authRouter.Get("/:id", handler.GetFacultyByID)
	authRouter.Post("/", handler.CreateFaculty)
	authRouter.Put("/:id", handler.UpdateFaculty)
}
