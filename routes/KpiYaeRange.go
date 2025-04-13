package routes

import (
	"byu-crm-service/modules/kpi-yae-range/http"
	"byu-crm-service/modules/kpi-yae-range/repository"
	"byu-crm-service/modules/kpi-yae-range/service"
	kpiYaeRepo "byu-crm-service/modules/kpi-yae/repository"
	kpiYaeService "byu-crm-service/modules/kpi-yae/service"
	visitHistoryRepo "byu-crm-service/modules/visit-history/repository"
	visitHistoryService "byu-crm-service/modules/visit-history/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func KpiYaeRangeRouter(router fiber.Router, db *gorm.DB) {
	kpiYaeRangeRepo := repository.NewKpiYaeRangeRepository(db)
	kpiYaeRepo := kpiYaeRepo.NewKpiYaeRepository(db)
	visitHistoryRepo := visitHistoryRepo.NewVisitHistoryRepository(db)
	kpiYaeRangeService := service.NewKpiYaeRangeService(kpiYaeRangeRepo)
	kpiYaeService := kpiYaeService.NewKpiYaeService(kpiYaeRepo)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)
	kpiYaeRangeHandler := http.NewKpiYaeRangeHandler(kpiYaeRangeService, kpiYaeService, visitHistoryService)

	http.KpiYaeRangeRoutes(router, kpiYaeRangeHandler)

}
