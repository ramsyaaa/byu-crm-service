package service

import (
	accountRepo "byu-crm-service/modules/account/repository"
	areaRepo "byu-crm-service/modules/area/repository"
	branchRepo "byu-crm-service/modules/branch/repository"
	cityRepo "byu-crm-service/modules/city/repository"
	clusterRepo "byu-crm-service/modules/cluster/repository"
	eligibilityRepo "byu-crm-service/modules/eligibility/repository"
	eligibilityResponse "byu-crm-service/modules/eligibility/response"
	"byu-crm-service/modules/product/repository"
	"byu-crm-service/modules/product/response"
	regionRepo "byu-crm-service/modules/region/repository"
	"fmt"
)

type productService struct {
	repo            repository.ProductRepository
	accountRepo     accountRepo.AccountRepository
	eligibilityRepo eligibilityRepo.EligibilityRepository
	areaRepo        areaRepo.AreaRepository
	regionRepo      regionRepo.RegionRepository
	branchRepo      branchRepo.BranchRepository
	clusterRepo     clusterRepo.ClusterRepository
	cityRepo        cityRepo.CityRepository
}

func NewProductService(
	repo repository.ProductRepository,
	accountRepo accountRepo.AccountRepository,
	eligibilityRepo eligibilityRepo.EligibilityRepository,
	areaRepo areaRepo.AreaRepository,
	regionRepo regionRepo.RegionRepository,
	branchRepo branchRepo.BranchRepository,
	clusterRepo clusterRepo.ClusterRepository,
	cityRepo cityRepo.CityRepository) ProductService {
	return &productService{
		repo:            repo,
		accountRepo:     accountRepo,
		eligibilityRepo: eligibilityRepo,
		areaRepo:        areaRepo,
		regionRepo:      regionRepo,
		branchRepo:      branchRepo,
		clusterRepo:     clusterRepo,
		cityRepo:        cityRepo,
	}
}

func (s *productService) GetAllProducts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, accountID int) ([]response.ProductResponse, int64, error) {
	var categories []string
	var types []string
	var locationFilter eligibilityResponse.LocationFilter

	if accountID != 0 {
		fmt.Println(accountID)
		account, err := s.accountRepo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))
		if err != nil {
			return nil, 0, err
		}
		if account.AccountCategory != nil {
			categories = append(categories, *account.AccountCategory)
		}

		if account.AccountType != nil {
			types = append(types, *account.AccountType)
		}

		if account.City != nil {
			city, err := s.cityRepo.GetCityByID(int(*account.City))
			if err != nil {
				return nil, 0, err
			}
			locationFilter = eligibilityResponse.LocationFilter{
				Cities: []string{city.Name},
			}
		}
	} else {
		if userRole == "Super-Admin" || userRole == "HQ" {
			locationFilter = eligibilityResponse.LocationFilter{}
		} else if userRole == "Area" {
			area, err := s.areaRepo.GetAreaByID(territoryID)
			if err != nil {
				return nil, 0, err
			}
			locationFilter = eligibilityResponse.LocationFilter{
				Areas: []string{area.Name},
			}
		} else if userRole == "Regional" {
			region, err := s.regionRepo.GetRegionByID(territoryID)
			if err != nil {
				return nil, 0, err
			}
			locationFilter = eligibilityResponse.LocationFilter{
				Regions: []string{region.Name},
			}
		} else if userRole == "Branch" || userRole == "Buddies" || userRole == "DS" || userRole == "YAE" || userRole == "Organic" {
			branch, err := s.branchRepo.GetBranchByID(territoryID)
			if err != nil {
				return nil, 0, err
			}
			locationFilter = eligibilityResponse.LocationFilter{
				Branches: []string{branch.Name},
			}
		} else if userRole == "Cluster" || userRole == "Admin-Tap" {
			cluster, err := s.clusterRepo.GetClusterByID(territoryID)
			if err != nil {
				return nil, 0, err
			}
			locationFilter = eligibilityResponse.LocationFilter{
				Clusters: []string{cluster.Name},
			}
		}
	}

	eligibilities, err := s.eligibilityRepo.GetEligibilities("App\\Models\\Product", categories, types, locationFilter)

	if err != nil {
		return nil, 0, err
	}

	var subjectIDs []int
	for _, e := range eligibilities {
		subjectIDs = append(subjectIDs, e.SubjectID)
	}

	return s.repo.GetAllProducts(limit, paginate, page, filters, subjectIDs)
}
