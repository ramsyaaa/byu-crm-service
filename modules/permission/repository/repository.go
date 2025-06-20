package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/permission/response"
)

type PermissionRepository interface {
	GetAllPermissions(limit int, paginate bool, page int, filters map[string]string) ([]response.PermissionResponse, int64, error)
	GetAllPermissionsByRoleID(role_id int) ([]response.PermissionResponse, error)
	GetPermissionByID(id int) (*response.PermissionResponse, error)
	GetPermissionByName(name string) (*response.PermissionResponse, error)
	CreatePermission(permission *models.Permission) (*response.PermissionResponse, error)
	UpdatePermission(oermission *models.Permission, id int) (*response.PermissionResponse, error)
	UpdateRolePermissions(roleID int, permissionIDs []int) error
	AddRolePermissions(roleID int, permissionIDs []int) error
}
