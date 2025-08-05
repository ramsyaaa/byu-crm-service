package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/branch/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) GetAllBranches(filters map[string]string, userRole string, territoryID int, withGeo bool) ([]response.BranchResponse, int64, error) {
	var branches []response.BranchResponse
	var total int64

	query := r.db.Model(&models.Branch{}).Joins("JOIN regions ON regions.id = branches.region_id")

	if withGeo {
		query = query.Select("regions.id, regions.name, regions.geojson")
	} else {
		query = query.Select("regions.id, regions.name")
	}

	// Filter berdasarkan role dan territory
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Cluster", "Admin-Tap":
			// Dari cluster → branch
			query = query.Joins("JOIN clusters ON clusters.branch_id = branches.id").
				Where("clusters.id = ?", territoryID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			// Langsung ke branch
			query = query.Where("branches.id = ?", territoryID)

		case "Region":
			// Dari region → branch
			query = query.Where("branches.region_id = ?", territoryID)

		case "Area":
			// Dari area → region → branch
			query = query.Where("regions.area_id = ?", territoryID)
		}
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("branches.name LIKE ?", "%"+token+"%")
		}
	}

	// Get total count before pagination
	query.Count(&total)

	err := query.Find(&branches).Error
	return branches, total, err
}

func (r *branchRepository) GetBranchByID(id int) (*response.BranchResponse, error) {
	var branch models.Branch
	err := r.db.First(&branch, id).Error
	if err != nil {
		return nil, err
	}

	branchResponse := &response.BranchResponse{
		ID:   branch.ID,
		Name: branch.Name,
		RegionID: func() int {
			if branch.RegionID != nil {
				region_id, _ := strconv.Atoi(*branch.RegionID)
				return region_id
			}
			return 0
		}(),
	}

	return branchResponse, nil
}

func (r *branchRepository) GetBranchByName(name string) (*response.BranchResponse, error) {
	var branch models.Branch
	err := r.db.Where("name = ?", name).First(&branch).Error
	if err != nil {
		return nil, err
	}

	branchResponse := &response.BranchResponse{
		ID:   branch.ID,
		Name: branch.Name,
		RegionID: func() int {
			if branch.RegionID != nil {
				regionID, _ := strconv.Atoi(*branch.RegionID)
				return regionID
			}
			return 0
		}(),
	}

	return branchResponse, nil
}

func (r *branchRepository) CreateBranch(branch *models.Branch) (*response.BranchResponse, error) {
	if err := r.db.Create(branch).Error; err != nil {
		return nil, err
	}

	var createdBranch models.Branch
	if err := r.db.First(&createdBranch, "id = ?", branch.ID).Error; err != nil {
		return nil, err
	}

	branchResponse := &response.BranchResponse{
		ID:   createdBranch.ID,
		Name: createdBranch.Name,
		RegionID: func() int {
			if createdBranch.RegionID != nil {
				regionID, _ := strconv.Atoi(*createdBranch.RegionID)
				return regionID
			}
			return 0
		}(),
	}

	return branchResponse, nil
}

func (r *branchRepository) UpdateBranch(branch *models.Branch, id int) (*response.BranchResponse, error) {
	var existingBranch models.Branch

	if err := r.db.First(&existingBranch, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingBranch).Updates(branch).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingBranch, "id = ?", id).Error; err != nil {
		return nil, err
	}

	branchResponse := &response.BranchResponse{
		ID:   existingBranch.ID,
		Name: existingBranch.Name,
		RegionID: func() int {
			if existingBranch.RegionID != nil {
				regionID, _ := strconv.Atoi(*existingBranch.RegionID)
				return regionID
			}
			return 0
		}(),
	}

	return branchResponse, nil
}
