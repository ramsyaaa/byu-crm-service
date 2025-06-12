package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/role/response"
)

type RoleRepository interface {
	GetAllRoles(limit int, paginate bool, page int, filters map[string]string) ([]response.RoleResponse, int64, error)
	GetRoleByID(id int) (*response.RoleResponse, error)
	GetRoleByName(name string) (*response.RoleResponse, error)
	CreateRole(role *models.Role) (*response.RoleResponse, error)
	UpdateRole(role *models.Role, id int) (*response.RoleResponse, error)
}
