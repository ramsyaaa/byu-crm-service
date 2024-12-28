package service

import "byu-crm-service/models"

type CityService interface {
	GetCityByID(id uint) (*models.City, error)
	GetCityByName(name string) (*models.City, error)
}
