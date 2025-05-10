package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/subdistrict/response"
)

type SubdistrictRepository interface {
	GetAllSubdistricts(filters map[string]string, userRole string, territoryID int) ([]response.SubdistrictResponse, int64, error)
	GetSubdistrictByID(id int) (*response.SubdistrictResponse, error)
	GetSubdistrictByName(name string) (*response.SubdistrictResponse, error)
	CreateSubdistrict(subdistrict *models.Subdistrict) (*response.SubdistrictResponse, error)
	UpdateSubdistrict(subdistrict *models.Subdistrict, id int) (*response.SubdistrictResponse, error)
}
