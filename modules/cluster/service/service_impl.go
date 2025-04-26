package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/cluster/repository"
	"byu-crm-service/modules/cluster/response"
	"fmt"
)

type clusterService struct {
	repo repository.ClusterRepository
}

func NewClusterService(repo repository.ClusterRepository) ClusterService {
	return &clusterService{repo: repo}
}

func (s *clusterService) GetAllClusters(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.ClusterResponse, int64, error) {
	return s.repo.GetAllClusters(limit, paginate, page, filters, userRole, territoryID)
}

func (s *clusterService) GetClusterByID(id int) (*response.ClusterResponse, error) {
	return s.repo.GetClusterByID(id)
}

func (s *clusterService) GetClusterByName(name string) (*response.ClusterResponse, error) {
	return s.repo.GetClusterByName(name)
}

func (s *clusterService) CreateCluster(name *string, branch_id int) (*response.ClusterResponse, error) {
	branchIDStr := fmt.Sprintf("%d", branch_id)
	cluster := &models.Cluster{Name: *name, BranchID: &branchIDStr}
	return s.repo.CreateCluster(cluster)
}

func (s *clusterService) UpdateCluster(name *string, branch_id int, id int) (*response.ClusterResponse, error) {
	branchIDStr := fmt.Sprintf("%d", branch_id)
	cluster := &models.Cluster{Name: *name, BranchID: &branchIDStr}
	return s.repo.UpdateCluster(cluster, id)
}
