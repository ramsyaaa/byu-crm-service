package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/city/response"
)

type CityRepository interface {
	GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.CityResponse, int64, error)
	GetCityByID(id int) (*response.CityResponse, error)
	GetCityByName(name string) (*response.CityResponse, error)
	CreateCity(city *models.City) (*response.CityResponse, error)
	UpdateCity(city *models.City, id int) (*response.CityResponse, error)
}
