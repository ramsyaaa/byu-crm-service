package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/territory/response"

	"gorm.io/gorm"
)

type territoryRepository struct {
	db *gorm.DB
}

func NewTerritoryRepository(db *gorm.DB) TerritoryRepository {
	return &territoryRepository{db: db}
}

func (r *territoryRepository) GetAllTerritories() (map[string]interface{}, int64, error) {
	var (
		areas    []response.AreaResponse
		regions  []response.RegionResponse
		branches []response.BranchResponse
		clusters []response.ClusterResponse
		cities   []response.CityResponse
		total    int64
	)

	if err := r.db.Model(&models.Area{}).Select("id", "name").Find(&areas).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&models.Region{}).Select("id", "name", "area_id").Find(&regions).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&models.Branch{}).Select("id", "name", "region_id").Find(&branches).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&models.Cluster{}).Select("id", "name", "branch_id").Find(&clusters).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&models.City{}).Select("id", "name", "cluster_id").Find(&cities).Error; err != nil {
		return nil, 0, err
	}

	result := map[string]interface{}{
		"areas":    areas,
		"regions":  regions,
		"branches": branches,
		"clusters": clusters,
		"cities":   cities,
	}

	return result, total, nil
}
