package service

import (
	"byu-crm-service/models"
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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

	var subjectIDs []int

	if userRole != "Super-Admin" && userRole != "HQ" {
		eligibilities, err := s.eligibilityRepo.GetEligibilities("App\\Models\\Product", categories, types, locationFilter)

		if err != nil {
			return nil, 0, err
		}

		for _, e := range eligibilities {
			subjectIDs = append(subjectIDs, e.SubjectID)
		}
	}

	return s.repo.GetAllProducts(limit, paginate, page, filters, subjectIDs)
}

func (s *productService) FindByProductID(id uint) (*response.SingleProductResponse, error) {
	product, err := s.repo.FindByProductID(id)
	if err != nil {
		return nil, err
	}

	var accountResponse response.SingleProductResponse

	accountResponse.ID = product.ID
	accountResponse.ProductName = *product.ProductName
	accountResponse.Description = *product.Description
	accountResponse.ProductCategory = *product.ProductCategory
	accountResponse.ProductType = *product.ProductType
	accountResponse.Bid = product.Bid
	accountResponse.Price = product.Price
	accountResponse.KeyVisual = product.KeyVisual
	accountResponse.AdditionalFile = product.AdditionalFile

	if product.QuotaValue != nil {
		quotaValueInt := int(*product.QuotaValue)
		quotaValueStr := fmt.Sprintf("%d", quotaValueInt)
		accountResponse.QuotaValue = &quotaValueStr
	}

	if product.ValidityValue != nil {
		validityValueInt := int(*product.ValidityValue)
		validityValueStr := fmt.Sprintf("%d", validityValueInt)
		accountResponse.ValidityValue = &validityValueStr
	}

	accountResponse.ValidityUnit = product.ValidityUnit
	accountResponse.CreatedAt = product.CreatedAt
	accountResponse.UpdatedAt = product.UpdatedAt

	accountResponse.EligibilityCategory = []string{}
	accountResponse.EligibilityType = []string{}
	accountResponse.EligibilityLocation = map[string][]string{}

	if product.Eligibility != nil {
		var categories []string
		if err := json.Unmarshal([]byte(product.Eligibility.Categories), &categories); err == nil {
			accountResponse.EligibilityCategory = categories
		}

		var types []string
		if err := json.Unmarshal([]byte(product.Eligibility.Types), &types); err == nil {
			accountResponse.EligibilityType = types
		}

		var locations map[string][]string
		if err := json.Unmarshal([]byte(product.Eligibility.Locations), &locations); err == nil {
			accountResponse.EligibilityLocation = locations
		}
	}

	return &accountResponse, nil
}

func (s *productService) CreateProduct(requestBody map[string]interface{}) (*models.Product, error) {
	// Use getStringValue to safely handle nil values and type conversions
	productData := map[string]string{
		"product_name":     getStringValue(requestBody["product_name"]),
		"bid":              getStringValue(requestBody["bid"]),
		"description":      getStringValue(requestBody["description"]),
		"price":            getStringValue(requestBody["price"]),
		"product_category": getStringValue(requestBody["product_category"]),
		"product_type":     getStringValue(requestBody["product_type"]),
		"key_visual":       getStringValue(requestBody["key_visual"]),
		"additional_file":  getStringValue(requestBody["additional_file"]),
		"quota_value":      getStringValue(requestBody["quota_value"]),
		"validity_value":   getStringValue(requestBody["validity_value"]),
		"validity_unit":    getStringValue(requestBody["validity_unit"]),
	}

	product, err := s.repo.CreateProduct(productData)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(requestBody map[string]interface{}, productID int) (*models.Product, error) {
	// Use getStringValue to safely handle nil values and type conversions
	productData := map[string]string{
		"product_name":     getStringValue(requestBody["product_name"]),
		"bid":              getStringValue(requestBody["bid"]),
		"description":      getStringValue(requestBody["description"]),
		"price":            getStringValue(requestBody["price"]),
		"product_category": getStringValue(requestBody["product_category"]),
		"product_type":     getStringValue(requestBody["product_type"]),
		"key_visual":       getStringValue(requestBody["key_visual"]),
		"additional_file":  getStringValue(requestBody["additional_file"]),
		"quota_value":      getStringValue(requestBody["quota_value"]),
		"validity_value":   getStringValue(requestBody["validity_value"]),
		"validity_unit":    getStringValue(requestBody["validity_unit"]),
	}

	product, err := s.repo.UpdateProduct(productData, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) InsertProductAccount(requestBody map[string]interface{}, account_id uint) ([]models.AccountProduct, error) {
	// Delete existing product accounts for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	productID, exists := requestBody["product_account"]
	if !exists {
		return nil, errors.New("product_account is missing")
	}

	var DataProductID []string

	switch v := productID.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataProductID = append(DataProductID, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataProductID
		DataProductID = append(DataProductID, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("product_id contains non-string value")
			}
			DataProductID = append(DataProductID, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid product_id type: %v", reflect.TypeOf(productID))
	}

	var productAccounts []models.AccountProduct

	// Loop through the contact IDs and create ProductAccount instances
	for _, contact := range DataProductID {
		idUint, err := strconv.ParseUint(contact, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting contact ID to uint: %v", err)
		}

		productAccounts = append(productAccounts, models.AccountProduct{
			ProductID: func(u uint) *uint { return &u }(uint(idUint)),
			AccountID: &account_id,
		})
	}

	// Insert the new contact accounts into the database
	if err := s.repo.Insert(productAccounts); err != nil {
		return nil, err
	}

	return productAccounts, nil
}

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
