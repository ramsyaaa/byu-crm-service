package routes

import (
	"byu-crm-service/modules/template-file/http"
	"byu-crm-service/modules/template-file/repository"
	"byu-crm-service/modules/template-file/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TemplateFileRouter(router fiber.Router, db *gorm.DB) {
	templateFileRepo := repository.NewTemplateFileRepository(db)

	templateFileService := service.NewTemplateFileService(templateFileRepo)

	templateFileHandler := http.NewTemplateFileHandler(templateFileService)

	http.TemplateFileRoutes(router, templateFileHandler)

}
