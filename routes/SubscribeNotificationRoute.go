package routes

import (
	"byu-crm-service/modules/notification-one-signal/http"
	"byu-crm-service/modules/notification-one-signal/repository"
	"byu-crm-service/modules/notification-one-signal/service"
	userRepo "byu-crm-service/modules/user/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubscribeNotificationRouter(router fiber.Router, db *gorm.DB) {
	notificationRepo := repository.NewNotificationOneSignalRepository(db)
	userRepository := userRepo.NewUserRepository(db)

	notificationService := service.NewNotificationOneSignalService(notificationRepo, userRepository)

	notificatioHandler := http.NewNotificationOneSignalHandler(notificationService)

	http.NotificationOneSignalRoutes(router, notificatioHandler)

}
