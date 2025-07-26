package routes

import (
	"byu-crm-service/modules/bak-generator/http"
	"byu-crm-service/modules/bak-generator/repository"
	"byu-crm-service/modules/bak-generator/service"
	cityRepo "byu-crm-service/modules/city/repository"

	accountRepo "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BakGeneratorRouter(router fiber.Router, db *gorm.DB) {
	bakRepo := repository.NewBakGeneratorRepository(db)
	bakService := service.NewBakGeneratorService(bakRepo)
	accountRepo := accountRepo.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	bakHandler := http.NewBakGeneratorHandler(bakService, accountService)

	http.BakGeneratorRoutes(router, bakHandler)

}
