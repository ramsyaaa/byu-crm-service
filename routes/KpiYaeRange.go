package routes

import (
	"byu-crm-service/modules/kpi-yae-range/http"
	"byu-crm-service/modules/kpi-yae-range/repository"
	"byu-crm-service/modules/kpi-yae-range/service"
	kpiYaeRepo "byu-crm-service/modules/kpi-yae/repository"
	kpiYaeService "byu-crm-service/modules/kpi-yae/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func KpiYaeRangeRouter(router fiber.Router, db *gorm.DB) {
	kpiYaeRangeRepo := repository.NewKpiYaeRangeRepository(db)
	kpiYaeRepo := kpiYaeRepo.NewKpiYaeRepository(db)
	kpiYaeRangeService := service.NewKpiYaeRangeService(kpiYaeRangeRepo)
	kpiYaeService := kpiYaeService.NewKpiYaeService(kpiYaeRepo)
	kpiYaeRangeHandler := http.NewKpiYaeRangeHandler(kpiYaeRangeService, kpiYaeService)

	http.KpiYaeRangeRoutes(router, kpiYaeRangeHandler)

}
