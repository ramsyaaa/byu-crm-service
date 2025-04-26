package repository

import (
	"byu-crm-service/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.City, int64, error) {
	var cities []models.City
	var total int64

	query := r.db.Model(&models.City{}).
		Joins("JOIN clusters ON clusters.id = cities.cluster_id").
		Joins("JOIN branches ON branches.id = clusters.branch_id").
		Joins("JOIN regions ON regions.id = branches.region_id").
		Joins("JOIN areas ON areas.id = regions.area_id")

	// Filter berdasarkan role dan territory ID jika tersedia
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Area":
			query = query.Where("areas.id = ?", territoryID)
		case "Region":
			query = query.Where("regions.id = ?", territoryID)
		case "Branch", "YAE", "Organic", "Buddies", "DS":
			query = query.Where("branches.id = ?", territoryID)
		case "Cluster", "Admin-Tap":
			query = query.Where("clusters.id = ?", territoryID)
		}
	}

	// Filter pencarian nama kota
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("cities.name LIKE ?", "%"+token+"%")
		}
	}

	// Filter tanggal
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("cities.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("cities.created_at <= ?", endDate)
	}

	// Hitung total sebelum paginasi
	query.Count(&total)

	// Sorting
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

	// Eksekusi query
	err := query.Find(&cities).Error
	return cities, total, err
}

func (r *cityRepository) FindByID(id uint) (*models.City, error) {
	var city models.City
	if err := r.db.First(&city, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &city, nil
}

func (r *cityRepository) FindByName(name string) (*models.City, error) {
	var city models.City
	if err := r.db.Where("name = ?", name).First(&city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &city, nil
}
