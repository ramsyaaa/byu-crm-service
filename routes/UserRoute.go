package routes

import (
	accountRepository "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"
	authRepository "byu-crm-service/modules/auth/repository"
	authService "byu-crm-service/modules/auth/service"
	cityRepo "byu-crm-service/modules/city/repository"
	"byu-crm-service/modules/user/http"
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func UserRouter(router fiber.Router, db *gorm.DB, redisClient any) {
	userRepo := repository.NewUserRepository(db)
	authRepo := authRepository.NewAuthRepository(db)
	accountRepo := accountRepository.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	userService := service.NewUserService(userRepo)
	authService := authService.NewAuthService(authRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	userHandler := http.NewUserHandler(userService, authService, accountService, redisClient.(*redis.Client))

	http.UserRoutes(router, userHandler)

}
