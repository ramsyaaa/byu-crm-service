package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/region/response"
)

type RegionRepository interface {
	GetAllRegions(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.RegionResponse, int64, error)
	GetRegionByID(id int) (*response.RegionResponse, error)
	GetRegionByName(name string) (*response.RegionResponse, error)
	CreateRegion(region *models.Region) (*response.RegionResponse, error)
	UpdateRegion(region *models.Region, id int) (*response.RegionResponse, error)
}
