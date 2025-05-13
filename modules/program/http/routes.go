package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProgramRoutes(router fiber.Router, handler *ProgramHandler) {
	authRouter := router.Group("/programs",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllPrograms)
	// authRouter.Get("/:id", handler.GetProgramById)
	// authRouter.Post("/", handler.CreateProgram)
	// authRouter.Put("/:id", handler.UpdateProgram)
}
