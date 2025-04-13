package routes

import (
	"byu-crm-service/modules/visit-checklist/http"
	"byu-crm-service/modules/visit-checklist/repository"
	"byu-crm-service/modules/visit-checklist/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func VisitChecklistRouter(router fiber.Router, db *gorm.DB) {
	visitChecklistRepo := repository.NewVisitChecklistRepository(db)

	visitChecklistService := service.NewVisitChecklistService(visitChecklistRepo)

	visitChecklistHandler := http.NewVisitChecklistHandler(visitChecklistService)

	http.VisitChecklistRoutes(router, visitChecklistHandler)

}
