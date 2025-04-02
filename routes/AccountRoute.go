package routes

import (
	"byu-crm-service/modules/account/http"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/service"
	cityRepo "byu-crm-service/modules/city/repository"
	contactAccountRepo "byu-crm-service/modules/contact-account/repository"
	socialMediaRepo "byu-crm-service/modules/social-media/repository"

	contactAccountService "byu-crm-service/modules/contact-account/service"
	socialMediaService "byu-crm-service/modules/social-media/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AccountRouter(router fiber.Router, db *gorm.DB) {
	accountRepo := repository.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	contactAccountRepo := contactAccountRepo.NewContactAccountRepository(db)
	socialMediaRepo := socialMediaRepo.NewSocialMediaRepository(db)

	accountService := service.NewAccountService(accountRepo, cityRepo)
	contactAccountService := contactAccountService.NewContactAccountService(contactAccountRepo)
	socialMediaService := socialMediaService.NewSocialMediaService(socialMediaRepo)

	accountHandler := http.NewAccountHandler(accountService, contactAccountService, socialMediaService)

	http.AccountRoutes(router, accountHandler)

}
