package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegistrationDealingRoutes(router fiber.Router, handler *RegistrationDealingHandler) {
	authRouter := router.Group(("/registration-dealings"),
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllRegistrationDealings)
	authRouter.Get("/:id", handler.GetRegistrationDealingById)
	authRouter.Post("/", handler.CreateRegistrationDealing)
}
