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

	absenceUserService := service.NewAbsenceUserService(absenceUserRepo)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	kpiYaeRangeService := kpiYaeRangeService.NewKpiYaeRangeService(kpiYaeRangeRepo)

	absenceUserHandler := http.NewAbsenceUserHandler(absenceUserService, visitHistoryService, accountService, kpiYaeRangeService)

	http.AbsenceUserRoutes(router, absenceUserHandler)

}
