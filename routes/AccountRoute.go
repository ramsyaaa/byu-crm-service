package routes

import (
	// "byu-crm-service/config" // Removed to fix import cycle
	absenceUserRepo "byu-crm-service/modules/absence-user/repository"
	accountFacultyRepo "byu-crm-service/modules/account-faculty/repository"
	accountMemberRepo "byu-crm-service/modules/account-member/repository"
	accountScheduleRepo "byu-crm-service/modules/account-schedule/repository"
	accountTypeCampusDetailRepo "byu-crm-service/modules/account-type-campus-detail/repository"
	accountTypeCommunityDetailRepo "byu-crm-service/modules/account-type-community-detail/repository"
	accountTypeSchoolDetailRepo "byu-crm-service/modules/account-type-school-detail/repository"
	"byu-crm-service/modules/account/http"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/service"
	"byu-crm-service/modules/account/validation"
	approvalLocationAccountRepo "byu-crm-service/modules/approval-location-account/repository"
	approvalLocationAccountService "byu-crm-service/modules/approval-location-account/service"
	areaRepo "byu-crm-service/modules/area/repository"
	branchRepo "byu-crm-service/modules/branch/repository"
	cityRepo "byu-crm-service/modules/city/repository"
	clusterRepo "byu-crm-service/modules/cluster/repository"
	contactAccountRepo "byu-crm-service/modules/contact-account/repository"
	eligibilityRepo "byu-crm-service/modules/eligibility/repository"
	notificationRepo "byu-crm-service/modules/notification/repository"
	notificationService "byu-crm-service/modules/notification/service"
	productRepo "byu-crm-service/modules/product/repository"
	regionRepo "byu-crm-service/modules/region/repository"
	smsSenderService "byu-crm-service/modules/sms-sender/service"
	socialMediaRepo "byu-crm-service/modules/social-media/repository"
	territoryRepo "byu-crm-service/modules/territory/repository"
	userRepo "byu-crm-service/modules/user/repository"

	absenceUserService "byu-crm-service/modules/absence-user/service"
	accountFacultyService "byu-crm-service/modules/account-faculty/service"
	accountMemberService "byu-crm-service/modules/account-member/service"
	accountScheduleService "byu-crm-service/modules/account-schedule/service"
	accountTypeCampusDetailService "byu-crm-service/modules/account-type-campus-detail/service"
	accountTypeCommunityDetailService "byu-crm-service/modules/account-type-community-detail/service"
	accountTypeSchoolDetailService "byu-crm-service/modules/account-type-school-detail/service"
	contactAccountService "byu-crm-service/modules/contact-account/service"
	productService "byu-crm-service/modules/product/service"
	socialMediaService "byu-crm-service/modules/social-media/service"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AccountRouter(router fiber.Router, db *gorm.DB, redisClient any) {
	// redisClient is now passed as a dependency to avoid import cycle

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
	productRepo := productRepo.NewProductRepository(db)
	eligibilityRepo := eligibilityRepo.NewEligibilityRepository(db)
	areaRepo := areaRepo.NewAreaRepository(db)
	regionRepo := regionRepo.NewRegionRepository(db)
	branchRepo := branchRepo.NewBranchRepository(db)
	clusterRepo := clusterRepo.NewClusterRepository(db)
	absenceUserRepo := absenceUserRepo.NewAbsenceUserRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	territoryRepo := territoryRepo.NewTerritoryRepository(db)
	approvalLocationAccountRepo := approvalLocationAccountRepo.NewApprovalLocationAccountRepository(db)

	// Set the account repository for validation
	validation.SetAccountRepository(accountRepo)

	accountService := service.NewAccountService(accountRepo, cityRepo)
	contactAccountService := contactAccountService.NewContactAccountService(contactAccountRepo)
	socialMediaService := socialMediaService.NewSocialMediaService(socialMediaRepo)
	accountTypeSchoolDetailService := accountTypeSchoolDetailService.NewAccountTypeSchoolDetailService(accountTypeSchoolDetailRepo)
	accountFacultyService := accountFacultyService.NewAccountFacultyService(accountFacultyRepo)
	accountMemberService := accountMemberService.NewAccountMemberService(accountMemberRepo)
	accountScheduleService := accountScheduleService.NewAccountScheduleService(accountScheduleRepo)
	accountTypeCampusDetailService := accountTypeCampusDetailService.NewAccountTypeCampusDetailService(accountTypeCampusDetailRepo)
	accountTypeCommunityDetailService := accountTypeCommunityDetailService.NewAccountTypeCommunityDetailService(accountTypeCommunityDetailRepo)
	productService := productService.NewProductService(productRepo, accountRepo, eligibilityRepo, areaRepo, regionRepo, branchRepo, clusterRepo, cityRepo)
	absenceUserService := absenceUserService.NewAbsenceUserService(absenceUserRepo, territoryRepo)
	userService := userService.NewUserService(userRepo)
	approvalLocationAccountService := approvalLocationAccountService.NewApprovalLocationAccountService(approvalLocationAccountRepo, accountRepo)

	accountHandler := http.NewAccountHandler(accountService, contactAccountService, socialMediaService, accountTypeSchoolDetailService, accountFacultyService, accountMemberService, accountScheduleService, accountTypeCampusDetailService, accountTypeCommunityDetailService, productService, absenceUserService, userService, approvalLocationAccountService, notificationService.NewNotificationService(notificationRepo.NewNotificationRepository(db), userRepo), smsSenderService.NewSmsSenderService(userRepo), redisClient.(*redis.Client))

	http.AccountRoutes(router, accountHandler)

}
