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

func (s *userService) GetUsersResume(only_role []string, userRole string, territoryID interface{}) (map[string]string, error) {
	users, _, err := s.repo.GetAllUsers(0, false, 1, map[string]string{}, only_role, false, userRole, territoryID)
	if err != nil {
		return nil, err
	}

	// Inisialisasi result dengan nilai awal "0" untuk setiap role di only_role
	result := make(map[string]string)
	for _, role := range only_role {
		result[role] = "0"
	}

	for _, user := range users {
		if len(user.RoleNames) > 0 {
			role := user.RoleNames[0]
			if _, ok := result[role]; ok {
				count, _ := strconv.Atoi(result[role])
				result[role] = strconv.Itoa(count + 1)
			}
		}
	}

	return result, nil
}

func (s *userService) UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error) {
	return s.repo.UpdateUserProfile(id, user)
}
