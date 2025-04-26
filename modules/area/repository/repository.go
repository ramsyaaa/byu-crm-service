package repository

import "byu-crm-service/models"

type AreaRepository interface {
	GetAllAreas(limit int, paginate bool, page int, filters map[string]string) ([]models.Area, int64, error)
	GetAreaByID(id int) (*models.Area, error)
	GetAreaByName(name string) (*models.Area, error)
	CreateArea(area *models.Area) (*models.Area, error)
	UpdateArea(area *models.Area, id int) (*models.Area, error)
}
