package service

import (
	"byu-crm-service/modules/city/response"
)

type CityService interface {
	GetAllCities(filters map[string]string, userRole string, territoryID int) ([]response.CityResponse, int64, error)
	GetCityByID(id int) (*response.CityResponse, error)
	GetCityByName(name string) (*response.CityResponse, error)
	CreateCity(name *string, cluster_id int) (*response.CityResponse, error)
	UpdateCity(name *string, cluster_id int, id int) (*response.CityResponse, error)
}
