package repository

import "byu-crm-service/modules/user/response"

type UserRepository interface {
	FindByID(id uint) (*response.UserResponse, error)
	GetAllUsers(limit int, paginate bool, page int, filters map[string]string, only_role []string, orderByMostAssignedPic bool, userRole string, territoryID interface{}) ([]response.UserResponse, int64, error)
	UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error)
	GetUserCountByRoles(onlyRoles []string, userRole string, territoryID interface{}) (map[string]int64, error)
}
