package http

import (
	"byu-crm-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ClusterRoutes(router fiber.Router, handler *ClusterHandler) {
	authRouter := router.Group("/clusters",
		middleware.JWTMiddleware,
		middleware.JWTUserContextMiddleware(),
	)
	authRouter.Get("/", handler.GetAllClusters)
	authRouter.Get("/:id", handler.GetClusterByID)
	authRouter.Post("/", handler.CreateCluster)
	authRouter.Put("/:id", handler.UpdateCluster)
}
