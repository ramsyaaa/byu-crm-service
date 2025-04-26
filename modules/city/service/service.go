package service

import "byu-crm-service/models"

type CityService interface {
	GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.City, int64, error)
	GetCityByID(id uint) (*models.City, error)
	GetCityByName(name string) (*models.City, error)
}
