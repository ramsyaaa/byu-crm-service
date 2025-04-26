package service

import (
	"byu-crm-service/modules/region/response"
)

type RegionService interface {
	GetAllRegions(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.RegionResponse, int64, error)
	GetRegionByID(id int) (*response.RegionResponse, error)
	GetRegionByName(name string) (*response.RegionResponse, error)
	CreateRegion(name *string, area_id int) (*response.RegionResponse, error)
	UpdateRegion(name *string, area_id int, id int) (*response.RegionResponse, error)
}
