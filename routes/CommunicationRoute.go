package routes

import (
	"byu-crm-service/modules/communication/http"
	"byu-crm-service/modules/communication/repository"
	"byu-crm-service/modules/communication/service"
	opportunityRepository "byu-crm-service/modules/opportunity/repository"
	opportunityService "byu-crm-service/modules/opportunity/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CommunicationRouter(router fiber.Router, db *gorm.DB) {
	communicationRepo := repository.NewCommunicationRepository(db)

	communicationService := service.NewCommunicationService(communicationRepo)
	opportunityRepo := opportunityRepository.NewOpportunityRepository(db)
	opportunityService := opportunityService.NewOpportunityService(opportunityRepo)

	communicationHandler := http.NewCommunicationHandler(communicationService, opportunityService)

	http.CommunicationRoutes(router, communicationHandler)

}
