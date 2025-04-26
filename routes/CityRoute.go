package routes

import (
	"byu-crm-service/modules/city/http"
	"byu-crm-service/modules/city/repository"
	"byu-crm-service/modules/city/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CityRouter(router fiber.Router, db *gorm.DB) {
	cityRepo := repository.NewCityRepository(db)
	cityService := service.NewCityService(cityRepo)
	cityHandler := http.NewCityHandler(cityService)

	http.CityRoutes(router, cityHandler)

}
