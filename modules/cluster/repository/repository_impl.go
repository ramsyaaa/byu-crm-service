package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/cluster/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type clusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db}
}

func (r *clusterRepository) GetAllClusters(filters map[string]string, userRole string, territoryID int, withGeo bool) ([]response.ClusterResponse, int64, error) {
	var clusters []response.ClusterResponse
	var total int64

	query := r.db.Model(&models.Cluster{}).
		Joins("JOIN branches ON branches.id = clusters.branch_id").
		Joins("JOIN regions ON regions.id = branches.region_id")

	if withGeo {
		query = query.Select("clusters.id, clusters.name, clusters.branch_id, clusters.geojson")
	} else {
		query = query.Select("clusters.id, clusters.name, clusters.branch_id")
	}

	// Filter berdasarkan role dan territory
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Cluster", "Admin-Tap":
			// Langsung dari cluster
			query = query.Where("clusters.id = ?", territoryID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			// Dari branch → cluster
			query = query.Where("clusters.branch_id = ?", territoryID)

		case "Region":
			// Dari region → branch → cluster
			query = query.Where("branches.region_id = ?", territoryID)

		case "Area":
			// Dari area → region → branch → cluster
			query = query.Where("regions.area_id = ?", territoryID)
		}
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("clusters.name LIKE ?", "%"+token+"%")
		}
	}

	// Get total count before pagination
	query.Count(&total)

	err := query.Find(&clusters).Error
	return clusters, total, err
}

func (r *clusterRepository) GetClusterByID(id int) (*response.ClusterResponse, error) {
	var cluster models.Cluster
	err := r.db.First(&cluster, id).Error
	if err != nil {
		return nil, err
	}

	clusterResponse := &response.ClusterResponse{
		ID:   cluster.ID,
		Name: cluster.Name,
		BranchID: func() int {
			if cluster.BranchID != nil {
				branchID, _ := strconv.Atoi(*cluster.BranchID)
				return branchID
			}
			return 0
		}(),
		Geojson: cluster.Geojson,
	}

	return clusterResponse, nil
}

func (r *clusterRepository) GetClusterByName(name string) (*response.ClusterResponse, error) {
	var cluster models.Cluster
	err := r.db.Where("name = ?", name).First(&cluster).Error
	if err != nil {
		return nil, err
	}

	clusterResponse := &response.ClusterResponse{
		ID:   cluster.ID,
		Name: cluster.Name,
		BranchID: func() int {
			if cluster.BranchID != nil {
				branchID, _ := strconv.Atoi(*cluster.BranchID)
				return branchID
			}
			return 0
		}(),
	}

	return clusterResponse, nil
}

func (r *clusterRepository) CreateCluster(cluster *models.Cluster) (*response.ClusterResponse, error) {
	if err := r.db.Create(cluster).Error; err != nil {
		return nil, err
	}

	var createdCluster models.Cluster
	if err := r.db.First(&createdCluster, "id = ?", cluster.ID).Error; err != nil {
		return nil, err
	}

	clusterResponse := &response.ClusterResponse{
		ID:   createdCluster.ID,
		Name: createdCluster.Name,
		BranchID: func() int {
			if createdCluster.BranchID != nil {
				branchID, _ := strconv.Atoi(*createdCluster.BranchID)
				return branchID
			}
			return 0
		}(),
	}

	return clusterResponse, nil
}

func (r *clusterRepository) UpdateCluster(cluster *models.Cluster, id int) (*response.ClusterResponse, error) {
	var existingCluster models.Cluster

	if err := r.db.First(&existingCluster, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingCluster).Updates(cluster).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingCluster, "id = ?", id).Error; err != nil {
		return nil, err
	}

	clusterResponse := &response.ClusterResponse{
		ID:   existingCluster.ID,
		Name: existingCluster.Name,
		BranchID: func() int {
			if existingCluster.BranchID != nil {
				branchID, _ := strconv.Atoi(*existingCluster.BranchID)
				return branchID
			}
			return 0
		}(),
	}

	return clusterResponse, nil
}
