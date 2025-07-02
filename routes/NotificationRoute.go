package routes

import (
	"byu-crm-service/modules/notification/http"
	"byu-crm-service/modules/notification/repository"
	"byu-crm-service/modules/notification/service"
	userRepo "byu-crm-service/modules/user/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NotificationRouter(router fiber.Router, db *gorm.DB) {
	notificatioRepo := repository.NewNotificationRepository(db)
	userRepo := userRepo.NewUserRepository(db)

	notificationService := service.NewNotificationService(notificatioRepo, userRepo)

	notificationHandler := http.NewNotificationHandler(notificationService)

	http.NotificationRoutes(router, notificationHandler)

}
