package routes

import (
	clusterRepo "byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/kpi-yae-range/http"
	"byu-crm-service/modules/kpi-yae-range/repository"
	"byu-crm-service/modules/kpi-yae-range/service"
	kpiYaeRepo "byu-crm-service/modules/kpi-yae/repository"
	kpiYaeService "byu-crm-service/modules/kpi-yae/service"
	performanceDigiposRepo "byu-crm-service/modules/performance-digipos/repository"
	performanceDigiposService "byu-crm-service/modules/performance-digipos/service"
	userRepo "byu-crm-service/modules/user/repository"
	userService "byu-crm-service/modules/user/service"
	visitHistoryRepo "byu-crm-service/modules/visit-history/repository"
	visitHistoryService "byu-crm-service/modules/visit-history/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func KpiYaeRangeRouter(router fiber.Router, db *gorm.DB) {
	kpiYaeRangeRepo := repository.NewKpiYaeRangeRepository(db)
	kpiYaeRepo := kpiYaeRepo.NewKpiYaeRepository(db)
	visitHistoryRepo := visitHistoryRepo.NewVisitHistoryRepository(db)
	performanceDigiposRepo := performanceDigiposRepo.NewPerformanceDigiposRepository(db)
	clusterRepo := clusterRepo.NewClusterRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	kpiYaeRangeService := service.NewKpiYaeRangeService(kpiYaeRangeRepo)
	kpiYaeService := kpiYaeService.NewKpiYaeService(kpiYaeRepo)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)
	performanceDigiposService := performanceDigiposService.NewPerformanceDigiposService(performanceDigiposRepo, clusterRepo)
	userService := userService.NewUserService(userRepo)
	kpiYaeRangeHandler := http.NewKpiYaeRangeHandler(kpiYaeRangeService, kpiYaeService, visitHistoryService, performanceDigiposService, userService)

	http.KpiYaeRangeRoutes(router, kpiYaeRangeHandler)

}
