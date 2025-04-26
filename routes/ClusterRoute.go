package routes

import (
	"byu-crm-service/modules/cluster/http"
	"byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/cluster/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClusterRouter(router fiber.Router, db *gorm.DB) {
	clusterRepo := repository.NewClusterRepository(db)

	clusterService := service.NewClusterService(clusterRepo)

	clusterHandler := http.NewClusterHandler(clusterService)

	http.ClusterRoutes(router, clusterHandler)

}
