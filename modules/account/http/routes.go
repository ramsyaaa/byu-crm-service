package http

import "github.com/gofiber/fiber/v2"

func AccountRoutes(app *fiber.App, handler *AccountHandler) {

	app.Post("/accounts/import", handler.Import)
	app.Get("/accounts", handler.GetAllAccounts)
}
