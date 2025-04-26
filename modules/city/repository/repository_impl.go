package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/city/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.CityResponse, int64, error) {
	var cities []response.CityResponse
	var total int64

	query := r.db.Model(&models.City{}).
		Joins("JOIN clusters ON clusters.id = cities.cluster_id").
		Joins("JOIN branches ON branches.id = clusters.branch_id").
		Joins("JOIN regions ON regions.id = branches.region_id")

	// Filter berdasarkan role dan territory
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Cluster", "Admin-Tap":
			query = query.Where("cities.cluster_id = ?", territoryID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			query = query.Where("clusters.branch_id = ?", territoryID)

		case "Region":
			query = query.Where("branches.region_id = ?", territoryID)

		case "Area":
			query = query.Where("regions.area_id = ?", territoryID)
		}
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("cities.name LIKE ?", "%"+token+"%")
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("cities.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("cities.created_at <= ?", endDate)
	}

	// Get total count before pagination
	query.Count(&total)

	// Ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy != "" && order != "" {
		query = query.Order(orderBy + " " + order)
	}

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&cities).Error
	return cities, total, err
}

func (r *cityRepository) GetCityByID(id int) (*response.CityResponse, error) {
	var city models.City
	err := r.db.First(&city, id).Error
	if err != nil {
		return nil, err
	}

	cityResponse := &response.CityResponse{
		ID:   city.ID,
		Name: city.Name,
		ClusterID: func() int {
			if city.ClusterID != nil {
				clusterID, _ := strconv.Atoi(*city.ClusterID)
				return clusterID
			}
			return 0
		}(),
	}

	return cityResponse, nil
}

func (r *cityRepository) GetCityByName(name string) (*response.CityResponse, error) {
	var city models.City
	err := r.db.Where("name = ?", name).First(&city).Error
	if err != nil {
		return nil, err
	}

	cityResponse := &response.CityResponse{
		ID:   city.ID,
		Name: city.Name,
		ClusterID: func() int {
			if city.ClusterID != nil {
				clusterID, _ := strconv.Atoi(*city.ClusterID)
				return clusterID
			}
			return 0
		}(),
	}

	return cityResponse, nil
}

func (r *cityRepository) CreateCity(city *models.City) (*response.CityResponse, error) {
	if err := r.db.Create(city).Error; err != nil {
		return nil, err
	}

	var createdCity models.City
	if err := r.db.First(&createdCity, "id = ?", city.ID).Error; err != nil {
		return nil, err
	}

	cityResponse := &response.CityResponse{
		ID:   createdCity.ID,
		Name: createdCity.Name,
		ClusterID: func() int {
			if createdCity.ClusterID != nil {
				clusterID, _ := strconv.Atoi(*createdCity.ClusterID)
				return clusterID
			}
			return 0
		}(),
	}

	return cityResponse, nil
}

func (r *cityRepository) UpdateCity(city *models.City, id int) (*response.CityResponse, error) {
	var existingCity models.City

	if err := r.db.First(&existingCity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingCity).Updates(city).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingCity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	cityResponse := &response.CityResponse{
		ID:   existingCity.ID,
		Name: existingCity.Name,
		ClusterID: func() int {
			if existingCity.ClusterID != nil {
				clusterID, _ := strconv.Atoi(*existingCity.ClusterID)
				return clusterID
			}
			return 0
		}(),
	}

	return cityResponse, nil
}
