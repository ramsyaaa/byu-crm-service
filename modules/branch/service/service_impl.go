package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/branch/repository"
	"byu-crm-service/modules/branch/response"
	"fmt"
)

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) GetAllBranches(filters map[string]string, userRole string, territoryID int) ([]response.BranchResponse, int64, error) {
	return s.repo.GetAllBranches(filters, userRole, territoryID)
}

func (s *branchService) GetBranchByID(id int) (*response.BranchResponse, error) {
	return s.repo.GetBranchByID(id)
}

func (s *branchService) GetBranchByName(name string) (*response.BranchResponse, error) {
	return s.repo.GetBranchByName(name)
}

func (s *branchService) CreateBranch(name *string, region_id int) (*response.BranchResponse, error) {
	regionIDStr := fmt.Sprintf("%d", region_id)
	branch := &models.Branch{Name: *name, RegionID: &regionIDStr}
	return s.repo.CreateBranch(branch)
}

func (s *branchService) UpdateBranch(name *string, region_id int, id int) (*response.BranchResponse, error) {
	regionIDStr := fmt.Sprintf("%d", region_id)
	branch := &models.Branch{Name: *name, RegionID: &regionIDStr}
	return s.repo.UpdateBranch(branch, id)
}
