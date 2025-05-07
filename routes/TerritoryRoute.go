package routes

import (
	"byu-crm-service/modules/territory/http"
	"byu-crm-service/modules/territory/repository"
	"byu-crm-service/modules/territory/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TerritoryRouter(router fiber.Router, db *gorm.DB) {
	territoryRepo := repository.NewTerritoryRepository(db)

	territoryService := service.NewTerritoryService(territoryRepo)

	territoryHandler := http.NewTerritoryHandler(territoryService)

	http.TerritoryRoutes(router, territoryHandler)

}
