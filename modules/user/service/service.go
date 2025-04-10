package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/repository"
)

type UserService interface {
	GetUserByID(id uint) (*repository.UserResponse, error)
	GetAllUsers(limit int, paginate bool, page int, filters map[string]string) ([]models.User, int64, error)
	UpdateUserProfile(id uint, user map[string]interface{}) (*repository.UserResponse, error)
}
