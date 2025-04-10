package repository

import (
	"byu-crm-service/models"
)

type UserRepository interface {
	FindByID(id uint) (*UserResponse, error)
	GetAllUsers(limit int, paginate bool, page int, filters map[string]string) ([]models.User, int64, error)
	UpdateUserProfile(id uint, user map[string]interface{}) (*UserResponse, error)
}
