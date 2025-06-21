package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/response"
)

type UserService interface {
	GetUserByID(id uint) (*response.UserResponse, error)
	GetAllUsers(limit int, paginate bool, page int, filters map[string]string, only_role []string, orderByMostAssignedPic bool, userRole string, territoryID interface{}) ([]*response.UserResponse, int64, error)
	UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error)
	GetUsersResume(only_role []string, userRole string, territoryID interface{}) (map[string]string, error)
	CreateUser(requestBody map[string]interface{}) (*models.User, error)
	UpdateUser(requestBody map[string]interface{}, userID int) (*response.UserResponse, error)
}
