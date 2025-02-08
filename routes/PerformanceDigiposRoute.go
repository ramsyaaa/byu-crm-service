package routes

import (
	clusterRepo "byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/performance-digipos/http"
	"byu-crm-service/modules/performance-digipos/repository"
	"byu-crm-service/modules/performance-digipos/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PerformanceDigiposRouter(app *fiber.App, db *gorm.DB) {
	performanceDigiposRepo := repository.NewPerformanceDigiposRepository(db)
	clusterRepo := clusterRepo.NewClusterRepository(db)
	performanceDigiposService := service.NewPerformanceDigiposService(performanceDigiposRepo, clusterRepo)
	performanceDigiposHandler := http.NewPerformanceDigiposHandler(performanceDigiposService)

	http.PerformanceDigiposRoutes(app, performanceDigiposHandler)

}
