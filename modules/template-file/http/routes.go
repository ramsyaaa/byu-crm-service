package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func TemplateFileRoutes(router fiber.Router, handler *TemplateFileHandler) {
	authRouter := router.Group("/template-file",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllFiles)
}
