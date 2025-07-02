package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"
	"byu-crm-service/modules/approval-location-account/http"
	"byu-crm-service/modules/approval-location-account/repository"
	"byu-crm-service/modules/approval-location-account/service"
	cityRepo "byu-crm-service/modules/city/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalLocationAccountRouter(router fiber.Router, db *gorm.DB) {
	approvalRepo := repository.NewApprovalLocationAccountRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)

	approvalService := service.NewApprovalLocationAccountService(approvalRepo, accountRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)

	approvalHandler := http.NewApprovalLocationAccountHandler(approvalService, accountService)

	http.ApprovalLocationAccountRoutes(router, approvalHandler)

}
