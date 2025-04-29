package routes

import (
	"byu-crm-service/modules/contact-account/http"
	"byu-crm-service/modules/contact-account/repository"
	"byu-crm-service/modules/contact-account/service"
	socialMediaRepo "byu-crm-service/modules/social-media/repository"
	socialMediaService "byu-crm-service/modules/social-media/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ContactRouter(router fiber.Router, db *gorm.DB) {
	cotactRepo := repository.NewContactAccountRepository(db)
	socialMediaRepo := socialMediaRepo.NewSocialMediaRepository(db)

	cotactService := service.NewContactAccountService(cotactRepo)
	socialMediaService := socialMediaService.NewSocialMediaService(socialMediaRepo)

	cotactHandler := http.NewContactAccountHandler(cotactService, socialMediaService)

	http.ContactAccountRoutes(router, cotactHandler)

}
