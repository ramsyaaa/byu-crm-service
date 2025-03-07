package http

import "github.com/gofiber/fiber/v2"

func DetailCommunityMemberRoute(app *fiber.App, handler *DetailCommunityMemberHandler) {

	app.Post("/detail-community-member/import", handler.Import)
}
