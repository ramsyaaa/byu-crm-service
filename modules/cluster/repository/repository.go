package repository

import "byu-crm-service/models"

type ClusterRepository interface {
	FindByID(id uint) (*models.Cluster, error)
	FindByName(name string) (*models.Cluster, error)
}
