package routes

import (
	"byu-crm-service/modules/absence-user/http"
	"byu-crm-service/modules/absence-user/repository"
	"byu-crm-service/modules/absence-user/service"
	visitHistoryRepo "byu-crm-service/modules/visit-history/repository"
	visitHistoryService "byu-crm-service/modules/visit-history/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AbsenceUserRouter(router fiber.Router, db *gorm.DB) {
	absenceUserRepo := repository.NewAbsenceUserRepository(db)
	visitHistoryRepo := visitHistoryRepo.NewVisitHistoryRepository(db)

	absenceUserService := service.NewAbsenceUserService(absenceUserRepo)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)

	absenceUserHandler := http.NewAbsenceUserHandler(absenceUserService, visitHistoryService)

	http.AbsenceUserRoutes(router, absenceUserHandler)

}
