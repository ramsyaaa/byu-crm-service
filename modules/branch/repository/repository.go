package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/branch/response"
)

type BranchRepository interface {
	GetAllBranches(filters map[string]string, userRole string, territoryID int) ([]response.BranchResponse, int64, error)
	GetBranchByID(id int) (*response.BranchResponse, error)
	GetBranchByName(name string) (*response.BranchResponse, error)
	CreateBranch(branch *models.Branch) (*response.BranchResponse, error)
	UpdateBranch(branch *models.Branch, id int) (*response.BranchResponse, error)
}
