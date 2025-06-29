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
	authRouter.Get("/overview", handler.GetCountAccount)
	authRouter.Get("/check-already-update-data/:id", handler.CheckAlreadyUpdateData)
	authRouter.Post("/import", handler.Import)
	authRouter.Get("/count-visited", handler.GetAccountVisitCounts)
	authRouter.Get("/", handler.GetAllAccounts)
	authRouter.Get("/:id", handler.GetAccountById)
	authRouter.Post("/update-pic/:id", handler.UpdatePic)
	authRouter.Post("/update-location/:id", handler.UpdateLocation)
	authRouter.Post("/update-pic-multiple", handler.UpdatePicMultipleAccounts)
	authRouter.Post("/update-priority", handler.UpdatePriorityMultipleAccounts)
	authRouter.Post("/", handler.CreateAccount)
	authRouter.Put("/:id", handler.UpdateAccount)
}
