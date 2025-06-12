package service

import (
	"byu-crm-service/modules/permission/response"
)

type PermissionService interface {
	GetAllPermissions(limit int, paginate bool, page int, filters map[string]string) ([]response.PermissionResponse, int64, error)
	GetAllPermissionsByRoleID(role_id int) ([]response.PermissionResponse, error)
	GetPermissionByID(id int) (*response.PermissionResponse, error)
	GetPermissionByName(name string) (*response.PermissionResponse, error)
	CreatePermission(name *string) (*response.PermissionResponse, error)
	UpdatePermission(name *string, id int) (*response.PermissionResponse, error)
	UpdateRolePermissions(roleID int, permissionIDs []int) error
	AddRolePermissions(roleID int, permissionIDs []int) error
}
