package repository

import "byu-crm-service/models"

type CityRepository interface {
	FindByID(id uint) (*models.City, error)
	FindByName(name string) (*models.City, error)
}
