package routes

import (
	"byu-crm-service/modules/role/http"
	"byu-crm-service/modules/role/repository"
	"byu-crm-service/modules/role/service"

	permissionRepo "byu-crm-service/modules/permission/repository"
	permissionService "byu-crm-service/modules/permission/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RoleRouter(router fiber.Router, db *gorm.DB, redisClient any) {
	roleRepository := repository.NewRoleRepository(db)
	permissionRepository := permissionRepo.NewPermissionRepository(db)

	roleService := service.NewRoleService(roleRepository)
	permissionService := permissionService.NewPermissionService(permissionRepository)

	RoleHandler := http.NewRoleHandler(roleService, permissionService, redisClient.(*redis.Client))

	http.RoleRoutes(router, RoleHandler)

}
