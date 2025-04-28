package routes

import (
	"byu-crm-service/modules/contact-account/http"
	"byu-crm-service/modules/contact-account/repository"
	"byu-crm-service/modules/contact-account/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ContactRouter(router fiber.Router, db *gorm.DB) {
	cotactRepo := repository.NewContactAccountRepository(db)

	cotactService := service.NewContactAccountService(cotactRepo)

	cotactHandler := http.NewContactAccountHandler(cotactService)

	http.ContactAccountRoutes(router, cotactHandler)

}
