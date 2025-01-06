package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/cluster/repository"
)

type clusterService struct {
	repo repository.ClusterRepository
}

func NewClusterService(repo repository.ClusterRepository) ClusterService {
	return &clusterService{repo: repo}
}

func (s *clusterService) GetClusterByID(id uint) (*models.Cluster, error) {
	return s.repo.FindByID(id)
}

func (s *clusterService) GetClusterByName(name string) (*models.Cluster, error) {
	return s.repo.FindByName(name)
}
