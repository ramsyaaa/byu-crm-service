package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SubdistrictRoutes(router fiber.Router, handler *SubdistrictHandler) {
	authRouter := router.Group("/subdistricts",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllSubdistricts)
	authRouter.Get("/:id", handler.GetSubdistrictByID)
	authRouter.Post("/", handler.CreateSubdistrict)
	authRouter.Put("/:id", handler.UpdateSubdistrict)
}
