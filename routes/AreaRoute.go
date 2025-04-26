package routes

import (
	"byu-crm-service/modules/area/http"
	"byu-crm-service/modules/area/repository"
	"byu-crm-service/modules/area/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AreaRouter(router fiber.Router, db *gorm.DB) {
	areaRepo := repository.NewAreaRepository(db)

	areaService := service.NewAreaService(areaRepo)

	areaHandler := http.NewAreaHandler(areaService)

	http.AreaRoutes(router, areaHandler)

}
