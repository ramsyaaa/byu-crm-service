package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/territory/http"
	"byu-crm-service/modules/territory/repository"
	"byu-crm-service/modules/territory/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TerritoryRouter(router fiber.Router, db *gorm.DB) {
	territoryRepo := repository.NewTerritoryRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)

	territoryService := service.NewTerritoryService(territoryRepo, accountRepo)

	territoryHandler := http.NewTerritoryHandler(territoryService)

	http.TerritoryRoutes(router, territoryHandler)

}
