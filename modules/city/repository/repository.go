package repository

import "byu-crm-service/models"

type CityRepository interface {
	GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.City, int64, error)
	FindByID(id uint) (*models.City, error)
	FindByName(name string) (*models.City, error)
}
