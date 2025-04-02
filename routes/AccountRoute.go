package routes

import (
	accountTypeSchoolDetailRepo "byu-crm-service/modules/account-type-school-detail/repository"
	"byu-crm-service/modules/account/http"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/service"
	cityRepo "byu-crm-service/modules/city/repository"
	contactAccountRepo "byu-crm-service/modules/contact-account/repository"
	socialMediaRepo "byu-crm-service/modules/social-media/repository"

	accountTypeSchoolDetailService "byu-crm-service/modules/account-type-school-detail/service"
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
	accountTypeSchoolDetailRepo := accountTypeSchoolDetailRepo.NewAccountTypeSchoolDetailRepository(db)

	accountService := service.NewAccountService(accountRepo, cityRepo)
	contactAccountService := contactAccountService.NewContactAccountService(contactAccountRepo)
	socialMediaService := socialMediaService.NewSocialMediaService(socialMediaRepo)
	accountTypeSchoolDetailService := accountTypeSchoolDetailService.NewAccountTypeSchoolDetailService(accountTypeSchoolDetailRepo)

	accountHandler := http.NewAccountHandler(accountService, contactAccountService, socialMediaService, accountTypeSchoolDetailService)

	http.AccountRoutes(router, accountHandler)

}
