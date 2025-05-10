package service

import (
	"byu-crm-service/modules/cluster/response"
)

type ClusterService interface {
	GetAllClusters(filters map[string]string, userRole string, territoryID int) ([]response.ClusterResponse, int64, error)
	GetClusterByID(id int) (*response.ClusterResponse, error)
	GetClusterByName(name string) (*response.ClusterResponse, error)
	CreateCluster(name *string, branch_id int) (*response.ClusterResponse, error)
	UpdateCluster(name *string, branch_id int, id int) (*response.ClusterResponse, error)
}
