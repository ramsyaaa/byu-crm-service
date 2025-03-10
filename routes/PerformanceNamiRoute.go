package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	cityRepo "byu-crm-service/modules/city/repository"
	"byu-crm-service/modules/performance-nami/http"
	"byu-crm-service/modules/performance-nami/repository"
	"byu-crm-service/modules/performance-nami/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PerformanceNamiRouter(app *fiber.App, db *gorm.DB) {
	performanceNamiRepo := repository.NewPerformanceNamiRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	performanceNamiService := service.NewPerformanceNamiService(performanceNamiRepo, cityRepo, accountRepo)
	performanceNamiHandler := http.NewPerformanceNamiHandler(performanceNamiService)

	http.PerformanceNamiRoutes(app, performanceNamiHandler)

}
