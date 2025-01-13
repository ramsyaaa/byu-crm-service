package service

import "byu-crm-service/models"

type SubdistrictService interface {
	GetSubdistrictByID(id uint) (*models.Subdistrict, error)
	GetSubdistrictByName(name string) (*models.Subdistrict, error)
}
