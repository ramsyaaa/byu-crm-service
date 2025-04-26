package service

import (
	"byu-crm-service/modules/subdistrict/response"
)

type SubdistrictService interface {
	GetAllSubdistricts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.SubdistrictResponse, int64, error)
	GetSubdistrictByID(id int) (*response.SubdistrictResponse, error)
	GetSubdistrictByName(name string) (*response.SubdistrictResponse, error)
	CreateSubdistrict(name *string, city_id int) (*response.SubdistrictResponse, error)
	UpdateSubdistrict(name *string, city_id int, id int) (*response.SubdistrictResponse, error)
}
