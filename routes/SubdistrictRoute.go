package routes

import (
	"byu-crm-service/modules/subdistrict/http"
	"byu-crm-service/modules/subdistrict/repository"
	"byu-crm-service/modules/subdistrict/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubdistrictRouter(router fiber.Router, db *gorm.DB) {
	subdistrictRepo := repository.NewSubdistrictRepository(db)
	subdistrictService := service.NewSubdistrictService(subdistrictRepo)
	subdistrictHandler := http.NewSubdistrictHandler(subdistrictService)

	http.SubdistrictRoutes(router, subdistrictHandler)

}
