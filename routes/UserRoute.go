package routes

import (
	authRepository "byu-crm-service/modules/auth/repository"
	authService "byu-crm-service/modules/auth/service"
	"byu-crm-service/modules/user/http"
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouter(router fiber.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	authRepo := authRepository.NewAuthRepository(db)
	userService := service.NewUserService(userRepo)
	authService := authService.NewAuthService(authRepo)
	userHandler := http.NewUserHandler(userService, authService)

	http.UserRoutes(router, userHandler)

}
