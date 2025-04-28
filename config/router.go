package config

import (
	"byu-crm-service/helper"
	"byu-crm-service/routes"
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
	})

	// Use the cors middleware to allow all origins and methods
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	if os.Getenv("APP_ENV") != "production" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Static("/static", "./static")

	// Serve the HTML dashboard on the root path
	app.Get("/log-viewer", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// Get available log files (uses helper function)
	app.Get("/logs", helper.GetLogFiles)

	// Get content from a specific log file (uses helper function)
	app.Get("/logs/:filename", helper.GetLogFileContent)

	api := fiber.New()
	api.Use(helper.LogToFile())
	// Register your routes here
	routes.PerformanceNamiRouter(api, db)
	routes.PerformanceSkulIdRouter(api, db)

	routes.PerformanceDigiposRouter(api, db)
	routes.DetailCommunityMemberRouter(api, db)

	routes.AuthRouter(api, db)

	routes.AreaRouter(api, db)
	routes.RegionRouter(api, db)
	routes.BranchRouter(api, db)
	routes.ClusterRouter(api, db)
	routes.CityRouter(api, db)
	routes.SubdistrictRouter(api, db)

	routes.AccountRouter(api, db)
	routes.FacultyRouter(api, db)
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
