package service

import (
	"byu-crm-service/modules/branch/response"
)

type BranchService interface {
	GetAllBranches(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.BranchResponse, int64, error)
	GetBranchByID(id int) (*response.BranchResponse, error)
	GetBranchByName(name string) (*response.BranchResponse, error)
	CreateBranch(name *string, region_id int) (*response.BranchResponse, error)
	UpdateBranch(name *string, region_id int, id int) (*response.BranchResponse, error)
}
