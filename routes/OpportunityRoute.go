package routes

import (
	"byu-crm-service/modules/opportunity/http"
	"byu-crm-service/modules/opportunity/repository"
	"byu-crm-service/modules/opportunity/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func OpportunityRouter(router fiber.Router, db *gorm.DB) {
	opportunityRepo := repository.NewOpportunityRepository(db)
	opportunityService := service.NewOpportunityService(opportunityRepo)
	opportunityHandler := http.NewOpportunityHandler(opportunityService)

	http.OpportunityRoutes(router, opportunityHandler)

}
