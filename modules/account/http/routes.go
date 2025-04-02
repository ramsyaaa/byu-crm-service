package http

import "github.com/gofiber/fiber/v2"

func AccountRoutes(router fiber.Router, handler *AccountHandler) {

	router.Post("/accounts/import", handler.Import)
	router.Get("/accounts", handler.GetAllAccounts)
	router.Post("/accounts", handler.CreateAccount)
}
