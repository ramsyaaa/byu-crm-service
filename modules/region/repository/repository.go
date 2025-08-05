package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/region/response"
)

type RegionRepository interface {
	GetAllRegions(filters map[string]string, userRole string, territoryID int, withGeo bool) ([]response.RegionResponse, int64, error)
	GetRegionByID(id int) (*response.RegionResponse, error)
	GetRegionByName(name string) (*response.RegionResponse, error)
	CreateRegion(region *models.Region) (*response.RegionResponse, error)
	UpdateRegion(region *models.Region, id int) (*response.RegionResponse, error)
}
