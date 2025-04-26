package routes

import (
	"byu-crm-service/modules/region/http"
	"byu-crm-service/modules/region/repository"
	"byu-crm-service/modules/region/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegionRouter(router fiber.Router, db *gorm.DB) {
	regionRepo := repository.NewRegionRepository(db)

	regionService := service.NewRegionService(regionRepo)

	regionHandler := http.NewRegionHandler(regionService)

	http.RegionRoutes(router, regionHandler)

}
