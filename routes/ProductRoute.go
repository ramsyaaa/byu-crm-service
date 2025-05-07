package routes

import (
	accountRepo "byu-crm-service/modules/account/repository"
	areaRepo "byu-crm-service/modules/area/repository"
	branchRepo "byu-crm-service/modules/branch/repository"
	cityRepo "byu-crm-service/modules/city/repository"
	clusterRepo "byu-crm-service/modules/cluster/repository"
	eligibilityRepo "byu-crm-service/modules/eligibility/repository"
	"byu-crm-service/modules/product/http"
	"byu-crm-service/modules/product/repository"
	"byu-crm-service/modules/product/service"
	regionRepo "byu-crm-service/modules/region/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProductRouter(router fiber.Router, db *gorm.DB) {
	productRepo := repository.NewProductRepository(db)
	accountRepo := accountRepo.NewAccountRepository(db)
	eligibilityRepo := eligibilityRepo.NewEligibilityRepository(db)
	areaRepo := areaRepo.NewAreaRepository(db)
	regionRepo := regionRepo.NewRegionRepository(db)
	branchRepo := branchRepo.NewBranchRepository(db)
	clusterRepo := clusterRepo.NewClusterRepository(db)
	cityRepo := cityRepo.NewCityRepository(db)

	productService := service.NewProductService(productRepo, accountRepo, eligibilityRepo, areaRepo, regionRepo, branchRepo, clusterRepo, cityRepo)

	productHandler := http.NewProductHandler(productService)

	http.ProductRoutes(router, productHandler)

}
