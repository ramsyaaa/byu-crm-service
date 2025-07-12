package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"
	"byu-crm-service/modules/approval-location-account/http"
	"byu-crm-service/modules/approval-location-account/repository"
	"byu-crm-service/modules/approval-location-account/service"
	cityRepo "byu-crm-service/modules/city/repository"
	notificationRepo "byu-crm-service/modules/notification/repository"
	notificationService "byu-crm-service/modules/notification/service"
	smsSenderService "byu-crm-service/modules/sms-sender/service"
	userRepo "byu-crm-service/modules/user/repository"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalLocationAccountRouter(router fiber.Router, db *gorm.DB) {
	approvalRepo := repository.NewApprovalLocationAccountRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	notificationRepo := notificationRepo.NewNotificationRepository(db)
	userRepo := userRepo.NewUserRepository(db)

	approvalService := service.NewApprovalLocationAccountService(approvalRepo, accountRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	notificationService := notificationService.NewNotificationService(notificationRepo, userRepo)
	smsSenderService := smsSenderService.NewSmsSenderService(userRepo)
	userService := userService.NewUserService(userRepo)

	approvalHandler := http.NewApprovalLocationAccountHandler(approvalService, accountService, notificationService, smsSenderService, userService)

	http.ApprovalLocationAccountRoutes(router, approvalHandler)

}
