package routes

import (
	"byu-crm-service/modules/yae-leaderboard/http"
	"byu-crm-service/modules/yae-leaderboard/repository"
	"byu-crm-service/modules/yae-leaderboard/service"

	userRepo "byu-crm-service/modules/user/repository"
	userService "byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func YaeLeaderboardRouter(router fiber.Router, db *gorm.DB) {
	yaeLeaderboardRepo := repository.NewYaeLeaderboardRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)

	yaeLeaderboardService := service.NewYaeLeaderboardService(yaeLeaderboardRepo)

	yaeLeaderboardHandler := http.NewYaeLeaderboardHandler(yaeLeaderboardService, userService)

	http.YaeLeaderboardRoutes(router, yaeLeaderboardHandler)

}
