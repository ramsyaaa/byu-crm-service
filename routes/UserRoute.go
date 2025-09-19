package routes

import (
	accountRepository "byu-crm-service/modules/account/repository"
	accountService "byu-crm-service/modules/account/service"
	authRepository "byu-crm-service/modules/auth/repository"
	authService "byu-crm-service/modules/auth/service"
	cityRepo "byu-crm-service/modules/city/repository"
	kpiYaeRangeRepo "byu-crm-service/modules/kpi-yae-range/repository"
	kpiYaeRangeService "byu-crm-service/modules/kpi-yae-range/service"
	roleRepo "byu-crm-service/modules/role/repository"
	roleService "byu-crm-service/modules/role/service"
	"byu-crm-service/modules/user/http"
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/service"
	"byu-crm-service/modules/user/validation"
	visitHistoryRepo "byu-crm-service/modules/visit-history/repository"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	yaeLeaderboardRepo "byu-crm-service/modules/yae-leaderboard/repository"
	yaeLeaderboardService "byu-crm-service/modules/yae-leaderboard/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func UserRouter(router fiber.Router, db *gorm.DB, redisClient any) {
	userRepo := repository.NewUserRepository(db)
	authRepo := authRepository.NewAuthRepository(db)
	accountRepo := accountRepository.NewAccountRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)
	roleRepo := roleRepo.NewRoleRepository(db)
	userService := service.NewUserService(userRepo)
	authService := authService.NewAuthService(authRepo)
	accountService := accountService.NewAccountService(accountRepo, cityRepo)
	roleService := roleService.NewRoleService(roleRepo)
	kpiYaeRangeRepo := kpiYaeRangeRepo.NewKpiYaeRangeRepository(db)
	kpiYaeRangeService := kpiYaeRangeService.NewKpiYaeRangeService(kpiYaeRangeRepo)
	visitHistoryRepo := visitHistoryRepo.NewVisitHistoryRepository(db)
	visitHistoryService := visitHistoryService.NewVisitHistoryService(visitHistoryRepo)
	yaeLeaderboardRepo := yaeLeaderboardRepo.NewYaeLeaderboardRepository(db)
	yaeLeaderboardService := yaeLeaderboardService.NewYaeLeaderboardService(yaeLeaderboardRepo)

	validation.SetUserRepository(userRepo)

	userHandler := http.NewUserHandler(userService, authService, accountService, roleService, kpiYaeRangeService, visitHistoryService, yaeLeaderboardService, redisClient.(*redis.Client))

	http.UserRoutes(router, userHandler)

}
