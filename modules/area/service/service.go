package service

import (
	"byu-crm-service/modules/area/response"
)

type AreaService interface {
	GetAllAreas(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.AreaResponse, int64, error)
	GetAreaByID(id int) (*response.AreaResponse, error)
	GetAreaByName(name string) (*response.AreaResponse, error)
	CreateArea(name *string) (*response.AreaResponse, error)
	UpdateArea(name *string, id int) (*response.AreaResponse, error)
}
