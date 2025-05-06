package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegistrationDealingRoutes(router fiber.Router, handler *RegistrationDealingHandler) {
	router.Post("/registration-dealings/no-auth", handler.CreateRegistrationDealing)

	authRouter := router.Group(("/registration-dealings"),
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllRegistrationDealings)
	authRouter.Get("/grouped", handler.GetAllRegistrationDealingsGrouped)
	authRouter.Get("/:id", handler.GetRegistrationDealingById)
	authRouter.Post("/", handler.CreateRegistrationDealing)

}
