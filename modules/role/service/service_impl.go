package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/role/repository"
	"byu-crm-service/modules/role/response"
)

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAllRoles(limit int, paginate bool, page int, filters map[string]string) ([]response.RoleResponse, int64, error) {
	return s.repo.GetAllRoles(limit, paginate, page, filters)
}

func (s *roleService) GetRoleByID(id int) (*response.RoleResponse, error) {
	return s.repo.GetRoleByID(id)
}

func (s *roleService) GetRoleByName(name string) (*response.RoleResponse, error) {
	return s.repo.GetRoleByName(name)
}

func (s *roleService) CreateRole(name *string) (*response.RoleResponse, error) {
	role := &models.Role{Name: *name, GuardName: "web"}
	return s.repo.CreateRole(role)
}

func (s *roleService) UpdateRole(name *string, id int) (*response.RoleResponse, error) {
	role := &models.Role{Name: *name, GuardName: "web"}
	return s.repo.UpdateRole(role, id)
}
