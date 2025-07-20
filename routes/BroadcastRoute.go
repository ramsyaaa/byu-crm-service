package routes

import (
	"byu-crm-service/modules/broadcast/http"
	"byu-crm-service/modules/broadcast/repository"
	"byu-crm-service/modules/broadcast/service"
	notificationRepo "byu-crm-service/modules/notification/repository"
	notificationService "byu-crm-service/modules/notification/service"
	roleRepo "byu-crm-service/modules/role/repository"
	roleService "byu-crm-service/modules/role/service"
	smsSender "byu-crm-service/modules/sms-sender/service"
	userRepo "byu-crm-service/modules/user/repository"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BroadcastRouter(router fiber.Router, db *gorm.DB) {
	broadcastRepo := repository.NewBroadcastRepository(db)
	broadcastService := service.NewBroadcastService(broadcastRepo)

	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)

	notificationRepo := notificationRepo.NewNotificationRepository(db)
	notificationService := notificationService.NewNotificationService(notificationRepo, userRepo)

	roleRepo := roleRepo.NewRoleRepository(db)
	roleService := roleService.NewRoleService(roleRepo)

	smsSenderService := smsSender.NewSmsSenderService(userRepo)

	broadcastHandler := http.NewBroadcastHandler(broadcastService, notificationService, userService, roleService, smsSenderService)

	http.BroadcastRoutes(router, broadcastHandler)

}
