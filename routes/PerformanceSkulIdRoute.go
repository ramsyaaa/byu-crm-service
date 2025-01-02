package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/performance-skul-id/http"
	"byu-crm-service/modules/performance-skul-id/repository"
	"byu-crm-service/modules/performance-skul-id/service"
	subdistrictRepo "byu-crm-service/modules/subdistrict/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PerformanceSkulIdRouter(app *fiber.App, db *gorm.DB) {
	performanceSkulIdRepo := repository.NewPerformanceSkulIdRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	subdistrictRepo := subdistrictRepo.NewSubdistrictRepository(db)
	performanceSkulIdService := service.NewPerformanceSkulIdService(performanceSkulIdRepo, accountRepo, subdistrictRepo)
	performanceSkulIdHandler := http.NewPerformanceSkulIdHandler(performanceSkulIdService)

	http.PerformanceSkulIdRoutes(app, performanceSkulIdHandler)

}
