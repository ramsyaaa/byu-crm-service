package routes

import (
	"byu-crm-service/modules/category/http"
	"byu-crm-service/modules/category/repository"
	"byu-crm-service/modules/category/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryRouter(router fiber.Router, db *gorm.DB) {
	categoryRepo := repository.NewCategoryRepository(db)

	categoryService := service.NewCategoryService(categoryRepo)

	categoryHandler := http.NewCategoryHandler(categoryService)

	http.CategoryRoutes(router, categoryHandler)

}
