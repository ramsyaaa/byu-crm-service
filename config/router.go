package config

import (
	"byu-crm-service/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	}))

	// Register your routes here
	routes.PerformanceNamiRouter(app, db)
	routes.PerformanceSkulIdRouter(app, db)
	routes.AccountRouter(app, db)
	routes.CityRouter(app, db)
	routes.SubdistrictRouter(app, db)
	routes.PerformanceDigiposRouter(app, db)
	routes.DetailCommunityMemberRouter(app, db)

	log.Fatalln(app.Listen(":" + os.Getenv("PORT")))
}
