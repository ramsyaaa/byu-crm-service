package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/area/response"
)

type AreaRepository interface {
	GetAllAreas(filters map[string]string, userRole string, territoryID int) ([]response.AreaResponse, int64, error)
	GetAreaByID(id int) (*response.AreaResponse, error)
	GetAreaByName(name string) (*response.AreaResponse, error)
	CreateArea(area *models.Area) (*response.AreaResponse, error)
	UpdateArea(area *models.Area, id int) (*response.AreaResponse, error)
}
