package routes

import (
	"byu-crm-service/modules/account/http"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/service"
	cityRepo "byu-crm-service/modules/city/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AccountRouter(app *fiber.App, db *gorm.DB) {
	accountRepo := repository.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	accountService := service.NewAccountService(accountRepo, cityRepo)
	accountHandler := http.NewAccountHandler(accountService)

	http.AccountRoutes(app, accountHandler)

}
