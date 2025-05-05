package routes

import (
	"byu-crm-service/modules/registration-dealing/http"
	"byu-crm-service/modules/registration-dealing/repository"
	"byu-crm-service/modules/registration-dealing/service"
	"byu-crm-service/modules/registration-dealing/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegistrationDealingRouter(router fiber.Router, db *gorm.DB) {
	registrationDealingRepo := repository.NewRegistrationDealingRepository(db)

	registrationDealingService := service.NewRegistrationDealingService(registrationDealingRepo)
	validation.SetRegistrationDealingRepository(registrationDealingRepo)

	registrationDealingHandler := http.NewRegistrationDealingHandler(registrationDealingService)

	http.RegistrationDealingRoutes(router, registrationDealingHandler)

}
