package routes

import (
	"byu-crm-service/modules/type/http"
	"byu-crm-service/modules/type/repository"
	"byu-crm-service/modules/type/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TypeRouter(router fiber.Router, db *gorm.DB) {
	typeRepo := repository.NewTypeRepository(db)

	typeService := service.NewTypeService(typeRepo)

	typeHandler := http.NewTypeHandler(typeService)

	http.TypeRoutes(router, typeHandler)

}
