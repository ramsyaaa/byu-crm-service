package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id uint) (*repository.UserResponse, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetAllUsers(limit int, paginate bool, page int, filters map[string]string) ([]models.User, int64, error) {
	return s.repo.GetAllUsers(limit, paginate, page, filters)
}

func (s *userService) UpdateUserProfile(id uint, user map[string]interface{}) (*repository.UserResponse, error) {
	return s.repo.UpdateUserProfile(id, user)
}
