package service

import (
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/response"
	"strconv"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetAllUsers(limit int, paginate bool, page int, filters map[string]string, only_role []string, orderByMostAssignedPic bool, userRole string, territoryID interface{}) ([]*response.UserResponse, int64, error) {
	users, total, err := s.repo.GetAllUsers(limit, paginate, page, filters, only_role, orderByMostAssignedPic, userRole, territoryID)
	if err != nil {
		return nil, 0, err
	}

	var userPointers []*response.UserResponse
	for i := range users {
		userPointers = append(userPointers, &users[i])
	}

	return userPointers, total, nil
}

func (s *userService) GetUsersResume(onlyRoles []string, userRole string, territoryID interface{}) (map[string]string, error) {
	counts, err := s.repo.GetUserCountByRoles(onlyRoles, userRole, territoryID)
	if err != nil {
		return nil, err
	}

	// Inisialisasi result dengan "0" dulu
	result := make(map[string]string)
	for _, role := range onlyRoles {
		result[role] = "0"
	}

	for role, count := range counts {
		result[role] = strconv.Itoa(int(count))
	}

	return result, nil
}

func (s *userService) UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error) {
	return s.repo.UpdateUserProfile(id, user)
}
