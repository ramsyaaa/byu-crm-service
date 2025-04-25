package routes

import (
	"byu-crm-service/modules/constant-data/http"
	"byu-crm-service/modules/constant-data/repository"
	"byu-crm-service/modules/constant-data/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ConstantDataRouter(router fiber.Router, db *gorm.DB) {
	constantDataRepo := repository.NewConstantDataRepository(db)

	constantDataService := service.NewConstantDataService(constantDataRepo)

	constantDataHandler := http.NewConstantDataHandler(constantDataService)

	http.ConstantDataRoutes(router, constantDataHandler)

}
