package service

import (
	"byu-crm-service/modules/role/response"
)

type RoleService interface {
	GetAllRoles(limit int, paginate bool, page int, filters map[string]string) ([]response.RoleResponse, int64, error)
	GetRoleByID(id int) (*response.RoleResponse, error)
	GetRoleByName(name string) (*response.RoleResponse, error)
	CreateRole(name *string) (*response.RoleResponse, error)
	UpdateRole(name *string, id int) (*response.RoleResponse, error)
}
