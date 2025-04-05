package routes

import (
	"byu-crm-service/modules/auth/http"
	"byu-crm-service/modules/auth/repository"
	"byu-crm-service/modules/auth/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRouter(router fiber.Router, db *gorm.DB) {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := http.NewAuthHandler(authService)

	http.AuthRoutes(router, authHandler)
}
