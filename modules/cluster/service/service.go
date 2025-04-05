package service

import "byu-crm-service/models"

type ClusterService interface {
	GetClusterByID(id uint) (*models.Cluster, error)
	GetClusterByName(name string) (*models.Cluster, error)
}
