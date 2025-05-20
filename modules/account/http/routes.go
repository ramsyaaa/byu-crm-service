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
	authRouter.Get("/overview", middleware.PermissionMiddleware("view account"), handler.GetCountAccount)
	authRouter.Get("/check-already-update-data/:id", handler.CheckAlreadyUpdateData)
	authRouter.Post("/import", middleware.PermissionMiddleware("import account"), handler.Import)
	authRouter.Get("/count-visited", middleware.PermissionMiddleware("view account"), handler.GetAccountVisitCounts)
	authRouter.Get("/", middleware.PermissionMiddleware("view account"), handler.GetAllAccounts)
	authRouter.Get("/:id", middleware.PermissionMiddleware("view account"), handler.GetAccountById)
	authRouter.Post("/update-pic/:id", middleware.PermissionMiddleware("edit account"), handler.UpdatePic)
	authRouter.Post("/", middleware.PermissionMiddleware("add account"), handler.CreateAccount)
	authRouter.Put("/:id", middleware.PermissionMiddleware("edit account"), handler.UpdateAccount)
}
