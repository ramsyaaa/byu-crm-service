package config

import (
	"byu-crm-service/helper"
	"byu-crm-service/middleware"
	"byu-crm-service/routes"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func Route(db *gorm.DB) {

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50 MB
		// Enable strict parsing for multipart forms and URL-encoded forms
		StrictRouting: true,
		// Add error handler for errors
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			// Log the error
			helper.LogError(c, fmt.Sprintf("Error in request: %s", err.Error()))

			// Return JSON response with error message
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Use the cors middleware to allow all origins and methods
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length, Content-Type",
	}))

	// Add recover middleware to catch panics
	app.Use(middleware.RecoverMiddleware())

	if os.Getenv("APP_ENV") != "production" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Static("/static", "./static")
	app.Static("/public", "./public")

	// Serve the HTML dashboard on the root path
	app.Get("/log-viewer", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// Get available log files (uses helper function)
	app.Get("/logs", helper.GetLogFiles)

	// Get content from a specific log file (uses helper function)
	app.Get("/logs/:filename", helper.GetLogFileContent)

	api := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50 MB
		// Enable strict parsing for multipart forms and URL-encoded forms
		StrictRouting: true,
		// Add error handler for errors
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			// Log the error with more details
			helper.LogError(c, fmt.Sprintf("API Error: %s", err.Error()))

			// Return JSON response with error message
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Add CORS middleware to the API router as well
	api.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length, Content-Type",
	}))

	// Add recover middleware to catch panics
	api.Use(middleware.RecoverMiddleware())

	// Add logger middleware
	api.Use(helper.LogToFile())
	// Register your routes here
	routes.PerformanceNamiRouter(api, db)
	routes.PerformanceSkulIdRouter(api, db)

	routes.PerformanceDigiposRouter(api, db)
	routes.DetailCommunityMemberRouter(api, db)

	routes.AuthRouter(api, db)

	routes.TerritoryRouter(api, db)
	routes.AreaRouter(api, db)
	routes.RegionRouter(api, db)
	routes.BranchRouter(api, db)
	routes.ClusterRouter(api, db)
	routes.CityRouter(api, db)
	routes.SubdistrictRouter(api, db)

	routes.AccountRouter(api, db)
	routes.ContactRouter(api, db)
	routes.OpportunityRouter(api, db)
	routes.CommunicationRouter(api, db)
	routes.FacultyRouter(api, db)

	routes.ProductRouter(api, db)

	routes.RegistrationDealingRouter(api, db)

	routes.AbsenceUserRouter(api, db)
	routes.UserRouter(api, db)
	routes.KpiYaeRangeRouter(api, db)
	routes.VisitChecklistRouter(api, db)
	routes.CategoryRouter(api, db)
	routes.TypeRouter(api, db)
	routes.ConstantDataRouter(api, db)

	app.Mount("/api/v1", api)
	log.Fatalln(app.Listen(":" + os.Getenv("PORT")))
}
