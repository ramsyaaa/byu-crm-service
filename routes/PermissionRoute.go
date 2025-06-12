package routes

import (
	"byu-crm-service/modules/permission/http"
	"byu-crm-service/modules/permission/repository"
	"byu-crm-service/modules/permission/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PermissionRouter(router fiber.Router, db *gorm.DB, redisClient any) {
	permissionRepo := repository.NewPermissionRepository(db)

	permissionService := service.NewPermissionService(permissionRepo)

	permissionHandler := http.NewPermissionHandler(permissionService, redisClient.(*redis.Client))

	http.PermissionRoutes(router, permissionHandler)

}
