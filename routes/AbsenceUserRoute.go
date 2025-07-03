package routes

import (
	"byu-crm-service/modules/absence-user/http"
	"byu-crm-service/modules/absence-user/repository"
	"byu-crm-service/modules/absence-user/service"
	accountRepo "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"
	cityRepo "byu-crm-service/modules/city/repository"
	kpiYaeRangeRepo "byu-crm-service/modules/kpi-yae-range/repository"
	kpiYaeRangeService "byu-crm-service/modules/kpi-yae-range/service"
	notificationRepo "byu-crm-service/modules/notification/repository"
	notificationService "byu-crm-service/modules/notification/service"
	smsSenderService "byu-crm-service/modules/sms-sender/service"
	territoryRepo "byu-crm-service/modules/territory/repository"
	userRepo "byu-crm-service/modules/user/repository"
	visitChecklistRepo "byu-crm-service/modules/visit-checklist/repository"
	visitChecklistService "byu-crm-service/modules/visit-checklist/service"
	visitHistoryRepo "byu-crm-service/modules/visit-history/repository"
	visitHistoryService "byu-crm-service/modules/visit-history/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AbsenceUserRouter(router fiber.Router, db *gorm.DB) {
	absenceUserRepo := repository.NewAbsenceUserRepository(db)
	visitHistoryRepo := visitHistoryRepo.NewVisitHistoryRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	kpiYaeRangeRepo := kpiYaeRangeRepo.NewKpiYaeRangeRepository(db)
	visitChecklistRepo := visitChecklistRepo.NewVisitChecklistRepository(db)
	territoryRepo := territoryRepo.NewTerritoryRepository(db)
	notificationRepo := notificationRepo.NewNotificationRepository(db)
	userRepo := userRepo.NewUserRepository(db)

	absenceUserService := service.NewAbsenceUserService(absenceUserRepo, territoryRepo)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	kpiYaeRangeService := kpiYaeRangeService.NewKpiYaeRangeService(kpiYaeRangeRepo)
	visitChecklistService := visitChecklistService.NewVisitChecklistService(visitChecklistRepo)
	notificationService := notificationService.NewNotificationService(notificationRepo, userRepo)

	absenceUserHandler := http.NewAbsenceUserHandler(absenceUserService, visitHistoryService, accountService, kpiYaeRangeService, visitChecklistService, notificationService, smsSenderService.NewSmsSenderService(userRepo))

	http.AbsenceUserRoutes(router, absenceUserHandler)

}
