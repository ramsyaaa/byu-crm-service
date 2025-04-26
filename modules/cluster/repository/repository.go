package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/cluster/response"
)

type ClusterRepository interface {
	GetAllClusters(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.ClusterResponse, int64, error)
	GetClusterByID(id int) (*response.ClusterResponse, error)
	GetClusterByName(name string) (*response.ClusterResponse, error)
	CreateCluster(cluster *models.Cluster) (*response.ClusterResponse, error)
	UpdateCluster(cluster *models.Cluster, id int) (*response.ClusterResponse, error)
}
