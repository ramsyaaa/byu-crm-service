package routes

import (
	"byu-crm-service/modules/menu/http"
	"byu-crm-service/modules/menu/repository"
	"byu-crm-service/modules/menu/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MenuRouter(router fiber.Router, db *gorm.DB) {
	menuRepo := repository.NewMenuRepository(db)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := http.NewMenuHandler(menuService)

	http.MenuRoutes(router, menuHandler)

}
