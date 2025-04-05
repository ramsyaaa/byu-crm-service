package routes

import (
	accountFacultyRepo "byu-crm-service/modules/account-faculty/repository"
	accountMemberRepo "byu-crm-service/modules/account-member/repository"
	accountScheduleRepo "byu-crm-service/modules/account-schedule/repository"
	accountTypeCampusDetailRepo "byu-crm-service/modules/account-type-campus-detail/repository"
	accountTypeCommunityDetailRepo "byu-crm-service/modules/account-type-community-detail/repository"
	accountTypeSchoolDetailRepo "byu-crm-service/modules/account-type-school-detail/repository"
	"byu-crm-service/modules/account/http"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/service"
	cityRepo "byu-crm-service/modules/city/repository"
	contactAccountRepo "byu-crm-service/modules/contact-account/repository"
	socialMediaRepo "byu-crm-service/modules/social-media/repository"

	accountFacultyService "byu-crm-service/modules/account-faculty/service"
	accountMemberService "byu-crm-service/modules/account-member/service"
	accountScheduleService "byu-crm-service/modules/account-schedule/service"
	accountTypeCampusDetailService "byu-crm-service/modules/account-type-campus-detail/service"
	accountTypeCommunityDetailService "byu-crm-service/modules/account-type-community-detail/service"
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
	accountFacultyRepo := accountFacultyRepo.NewAccountFacultyRepository(db)
	accountMemberRepo := accountMemberRepo.NewAccountMemberRepository(db)
	accountScheduleRepo := accountScheduleRepo.NewAccountScheduleRepository(db)
	accountTypeCampusDetailRepo := accountTypeCampusDetailRepo.NewAccountTypeCampusDetailRepository(db)
	accountTypeCommunityDetailRepo := accountTypeCommunityDetailRepo.NewAccountTypeCommunityDetailRepository(db)

	accountService := service.NewAccountService(accountRepo, cityRepo)
	contactAccountService := contactAccountService.NewContactAccountService(contactAccountRepo)
	socialMediaService := socialMediaService.NewSocialMediaService(socialMediaRepo)
	accountTypeSchoolDetailService := accountTypeSchoolDetailService.NewAccountTypeSchoolDetailService(accountTypeSchoolDetailRepo)
	accountFacultyService := accountFacultyService.NewAccountFacultyService(accountFacultyRepo)
	accountMemberService := accountMemberService.NewAccountMemberService(accountMemberRepo)
	accountScheduleService := accountScheduleService.NewAccountScheduleService(accountScheduleRepo)
	accountTypeCampusDetailService := accountTypeCampusDetailService.NewAccountTypeCampusDetailService(accountTypeCampusDetailRepo)
	accountTypeCommunityDetailService := accountTypeCommunityDetailService.NewAccountTypeCommunityDetailService(accountTypeCommunityDetailRepo)

	accountHandler := http.NewAccountHandler(accountService, contactAccountService, socialMediaService, accountTypeSchoolDetailService, accountFacultyService, accountMemberService, accountScheduleService, accountTypeCampusDetailService, accountTypeCommunityDetailService)

	http.AccountRoutes(router, accountHandler)

}
