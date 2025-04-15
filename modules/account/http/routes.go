package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func AccountRoutes(router fiber.Router, handler *AccountHandler) {
	authRouter := router.Group("/accounts",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Post("/import", handler.Import)
	authRouter.Get("/count-visited", handler.GetAccountVisitCounts)
	authRouter.Get("/", handler.GetAllAccounts)
	authRouter.Get("/:id", handler.GetAccountById)
	authRouter.Post("/", handler.CreateAccount)
	authRouter.Put("/:id", handler.UpdateAccount)
}
