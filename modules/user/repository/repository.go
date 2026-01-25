package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/response"
)

type UserRepository interface {
	FindByID(id uint) (*response.UserResponse, error)
	GetUserByIDs(ids []uint) ([]response.UserResponse, error)
	FindByEmail(email string) (*response.UserResponse, error)
	FindByMsisdn(msisdn string) (*response.UserResponse, error)
	GetAllUsers(limit int, paginate bool, page int, filters map[string]string, only_role []string, orderByMostAssignedPic bool, userRole string, territoryID interface{}) ([]response.UserResponse, int64, error)
	CreateUser(requestBody map[string]string) (*models.User, error)
	UpdateUser(requestBody map[string]string, userID int) (*response.UserResponse, error)
	UpdateYaeCode(userID uint, yaeCode string) error
	UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error)
	GetUserCountByRoles(onlyRoles []string, userRole string, territoryID interface{}) (map[string]int64, error)
	ResignUser(id uint) error
}
