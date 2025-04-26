package routes

import (
	"byu-crm-service/modules/branch/http"
	"byu-crm-service/modules/branch/repository"
	"byu-crm-service/modules/branch/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BranchRouter(router fiber.Router, db *gorm.DB) {
	branchRepo := repository.NewBranchRepository(db)

	branchService := service.NewBranchService(branchRepo)

	branchHandler := http.NewBranchHandler(branchService)

	http.BranchRoutes(router, branchHandler)

}
