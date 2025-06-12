package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/permission/repository"
	"byu-crm-service/modules/permission/response"
)

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) GetAllPermissions(limit int, paginate bool, page int, filters map[string]string) ([]response.PermissionResponse, int64, error) {
	return s.repo.GetAllPermissions(limit, paginate, page, filters)
}

func (s *permissionService) GetAllPermissionsByRoleID(role_id int) ([]response.PermissionResponse, error) {
	return s.repo.GetAllPermissionsByRoleID(role_id)
}

func (s *permissionService) GetPermissionByID(id int) (*response.PermissionResponse, error) {
	return s.repo.GetPermissionByID(id)
}

func (s *permissionService) GetPermissionByName(name string) (*response.PermissionResponse, error) {
	return s.repo.GetPermissionByName(name)
}

func (s *permissionService) CreatePermission(name *string) (*response.PermissionResponse, error) {
	Permission := &models.Permission{Name: *name, GuardName: "web"}
	return s.repo.CreatePermission(Permission)
}

func (s *permissionService) UpdatePermission(name *string, id int) (*response.PermissionResponse, error) {
	Permission := &models.Permission{Name: *name, GuardName: "web"}
	return s.repo.UpdatePermission(Permission, id)
}

func (s *permissionService) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	return s.repo.UpdateRolePermissions(roleID, permissionIDs)
}
func (s *permissionService) AddRolePermissions(roleID int, permissionIDs []int) error {
	return s.repo.AddRolePermissions(roleID, permissionIDs)
}
