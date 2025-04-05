package routes

import (
	"byu-crm-service/modules/user/http"
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouter(router fiber.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	http.UserRoutes(router, userHandler)

}
