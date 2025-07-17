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

	RedisClient := InitRedis()

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50 MB
		// Disable strict routing to allow more flexible URL handling
		StrictRouting: false,
		// Add error handler for errors
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			// Log the error to database
			helper.LogErrorToDatabase(db, c, fmt.Sprintf("Error in request: %s", err.Error()))

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
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders: "Content-Length, Content-Type",
		MaxAge:        86400, // 24 hours
	}))

	// Add database recover middleware to catch panics
	app.Use(middleware.DatabaseRecoverMiddleware(db))

	if os.Getenv("APP_ENV") != "production" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Static("/static", "./static")
	app.Static("/public", "./public")

	// Root URL redirect to admin login
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/login")
	})

	// Admin routes with authentication
	adminGroup := app.Group("/admin")

	// Admin login page (no authentication required)
	adminGroup.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("./static/admin-login.html")
	})

	// Protected admin routes
	adminProtected := adminGroup.Group("", middleware.JWTMiddleware, middleware.AdminAuthMiddleware())

	// Admin dashboard (default landing page)
	adminProtected.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendFile("./static/admin-dashboard.html")
	})

	// Redirect /admin to /admin/dashboard
	adminGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/dashboard")
	})

	// Legacy log viewer route (redirect to admin dashboard)
	app.Get("/log-viewer", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/dashboard")
	})

	// Database log viewer endpoints
	logHandler := helper.NewLogViewerHandler(db)
	app.Get("/api-logs", logHandler.GetApiLogs)
	app.Get("/api-logs/stats", logHandler.GetLogStats)
	app.Get("/api-logs/errors", logHandler.GetErrorLogs)
	app.Get("/api-logs/slow", logHandler.GetSlowRequests)
	app.Get("/api-logs/:id", logHandler.GetLogById)
	app.Post("/api-logs/cleanup", logHandler.CleanupLogs)

	// Chart data endpoints
	app.Get("/api-logs/chart-data/requests-over-time", logHandler.GetRequestsOverTime)
	app.Get("/api-logs/chart-data/status-distribution", logHandler.GetStatusDistribution)

	// MAU (Monthly Active Users) endpoints
	app.Get("/api-logs/mau", logHandler.GetMAUData)
	app.Get("/api-logs/users", logHandler.GetUsersList)
	app.Get("/api-logs/user-activity", logHandler.GetUserActivityTimeline)

	api := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50 MB
		// Disable strict routing to allow more flexible URL handling
		StrictRouting: false,
		// Add error handler for errors
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			// Log the error to database with more details
			helper.LogErrorToDatabase(db, c, fmt.Sprintf("API Error: %s", err.Error()))

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
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders: "Content-Length, Content-Type",
		MaxAge:        86400, // 24 hours
	}))

	// Add database recover middleware to catch panics
	api.Use(middleware.DatabaseRecoverMiddleware(db))

	// Add database logger middleware
	api.Use(helper.DatabaseLogger(db))
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

	routes.AccountRouter(api, db, RedisClient)
	routes.ContactRouter(api, db)
	routes.OpportunityRouter(api, db)
	routes.CommunicationRouter(api, db)
	routes.FacultyRouter(api, db)

	routes.ProductRouter(api, db)

	routes.RegistrationDealingRouter(api, db)

	routes.AbsenceUserRouter(api, db)
	routes.UserRouter(api, db, RedisClient)
	routes.KpiYaeRangeRouter(api, db)
	routes.VisitChecklistRouter(api, db)
	routes.CategoryRouter(api, db)
	routes.TypeRouter(api, db)
	routes.RoleRouter(api, db, RedisClient)
	routes.PermissionRouter(api, db, RedisClient)
	routes.ConstantDataRouter(api, db)
	routes.MenuRouter(api, db)
	routes.ApprovalLocationAccountRouter(api, db)
	routes.NotificationRouter(api, db)

	app.Mount("/api/v1", api)
	log.Fatalln(app.Listen(":" + os.Getenv("PORT")))
}
