package repository

import "byu-crm-service/models"

type SubdistrictRepository interface {
	FindByID(id uint) (*models.Subdistrict, error)
	FindByName(name string) (*models.Subdistrict, error)
}
