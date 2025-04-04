package routes

import (
	"byu-crm-service/modules/absence-user/http"
	"byu-crm-service/modules/absence-user/repository"
	"byu-crm-service/modules/absence-user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AbsenceUserRouter(router fiber.Router, db *gorm.DB) {
	absenceUserRepo := repository.NewAbsenceUserRepository(db)

	absenceUserService := service.NewAbsenceUserService(absenceUserRepo)

	absenceUserHandler := http.NewAbsenceUserHandler(absenceUserService)

	http.AbsenceUserRoutes(router, absenceUserHandler)

}
