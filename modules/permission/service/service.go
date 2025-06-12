package service

import (
	"byu-crm-service/modules/permission/response"
)

type PermissionService interface {
	GetAllPermissions(limit int, paginate bool, page int, filters map[string]string) ([]response.PermissionResponse, int64, error)
	GetPermissionByID(id int) (*response.PermissionResponse, error)
	GetPermissionByName(name string) (*response.PermissionResponse, error)
	CreatePermission(name *string) (*response.PermissionResponse, error)
	UpdatePermission(name *string, id int) (*response.PermissionResponse, error)
}
