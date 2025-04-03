package routes

import (
	"byu-crm-service/modules/faculty/http"
	"byu-crm-service/modules/faculty/repository"
	"byu-crm-service/modules/faculty/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FacultyRouter(router fiber.Router, db *gorm.DB) {
	facultyRepo := repository.NewFacultyRepository(db)

	facultyService := service.NewFacultyService(facultyRepo)

	facultyHandler := http.NewFacultyHandler(facultyService)

	http.FacultyRoutes(router, facultyHandler)

}
