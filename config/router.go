package config

import (
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

	// Register your routes here
	routes.PerformanceNamiRouter(app, db)
	routes.PerformanceSkulIdRouter(app, db)

	routes.PerformanceDigiposRouter(app, db)
	routes.DetailCommunityMemberRouter(app, db)

	authGroup := app.Group("/api/v1")
	routes.AuthRouter(authGroup, db)

	routes.AreaRouter(authGroup, db)
	routes.RegionRouter(authGroup, db)
	routes.BranchRouter(authGroup, db)
	routes.ClusterRouter(authGroup, db)
	routes.CityRouter(authGroup, db)
	routes.SubdistrictRouter(authGroup, db)

	routes.AccountRouter(authGroup, db)
	routes.FacultyRouter(authGroup, db)
	routes.AbsenceUserRouter(authGroup, db)
	routes.UserRouter(authGroup, db)
	routes.KpiYaeRangeRouter(authGroup, db)
	routes.VisitChecklistRouter(authGroup, db)
	routes.CategoryRouter(authGroup, db)
	routes.TypeRouter(authGroup, db)
	routes.ConstantDataRouter(authGroup, db)

	log.Fatalln(app.Listen(":" + os.Getenv("PORT")))
}
