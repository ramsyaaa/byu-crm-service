package routes

import (
	"byu-crm-service/modules/performance-indiana/http"
	"byu-crm-service/modules/performance-indiana/repository"
	"byu-crm-service/modules/performance-indiana/service"
	userRepo "byu-crm-service/modules/user/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PerformanceIndianaRouter(app *fiber.App, db *gorm.DB) {
	performanceIndianaRepo := repository.NewPerformanceIndianaRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	performanceIndianaService := service.NewPerformanceIndianaService(performanceIndianaRepo, userRepo)
	performanceIndianaHandler := http.NewPerformanceIndianaHandler(performanceIndianaService)

	http.PerformanceIndianaRoutes(app, performanceIndianaHandler)

}
