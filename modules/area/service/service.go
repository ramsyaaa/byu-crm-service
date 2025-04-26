package service

import "byu-crm-service/models"

type AreaService interface {
	GetAllAreas(limit int, paginate bool, page int, filters map[string]string) ([]models.Area, int64, error)
	GetAreaByID(id int) (*models.Area, error)
	GetAreaByName(name string) (*models.Area, error)
	CreateArea(name *string) (*models.Area, error)
	UpdateArea(name *string, id int) (*models.Area, error)
}
