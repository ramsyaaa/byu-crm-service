package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AbsenceUserRoutes(router fiber.Router, handler *AbsenceUserHandler) {
	// Create a group for absence-user routes with CORS middleware
	absenceGroup := router.Group("/absence-user")

	// Add CORS middleware specifically for absence-user routes
	absenceGroup.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders: "Content-Length, Content-Type",
		MaxAge:        86400, // 24 hours
	}))

	// Add OPTIONS handler for preflight requests
	absenceGroup.Options("/*", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		return c.SendStatus(fiber.StatusOK)
	})

	// Create authenticated routes
	authRouter := absenceGroup.Group("/",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)

	// Register routes
	authRouter.Get("/", handler.GetAllAbsenceUsers)
	authRouter.Get("/active-absence", handler.GetAbsenceActive)
	authRouter.Get("/:id", handler.GetAbsenceUserByID)
	authRouter.Post("/", handler.CreateAbsenceUser)
}
