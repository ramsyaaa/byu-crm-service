package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/repository"
	"byu-crm-service/modules/user/response"
	"fmt"
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

func (s *userService) CreateUser(requestBody map[string]interface{}) (*models.User, error) {
	// Use getStringValue to safely handle nil values and type conversions
	var territory_type string
	var territory_id int
	if requestBody["role_id"] != "1" && requestBody["role_id"] != "2" && requestBody["role_id"] != "12" {
		// Role Area
		if requestBody["role_id"] == "3" {
			territory_type = "App\\Models\\Area"
			if v, ok := requestBody["area_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "4" {
			territory_type = "App\\Models\\Region"
			if v, ok := requestBody["region_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "5" || requestBody["role_id"] == "7" || requestBody["role_id"] == "9" || requestBody["role_id"] == "10" || requestBody["role_id"] == "11" {
			territory_type = "App\\Models\\Branch"
			if v, ok := requestBody["branch_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "6" || requestBody["role_id"] == "8" {
			territory_type = "App\\Models\\Cluster"
			if v, ok := requestBody["cluster_id"].(float64); ok {
				territory_id = int(v)
			}
		}
	}
	userData := map[string]string{
		"name":              getStringValue(requestBody["name"]),
		"email":             getStringValue(requestBody["email"]),
		"msisdn":            getStringValue(requestBody["msisdn"]),
		"user_status":       "active",
		"user_type":         getStringValue(requestBody["user_type"]),
		"territory_type":    territory_type,
		"territory_id":      getStringValue(territory_id),
		"outlet_id_digipos": getStringValue(requestBody["outlet_id_digipos"]),
		"nami_agent_id":     getStringValue(requestBody["nami_agent_id"]),
		"password":          getStringValue(requestBody["password"]),
	}

	accounts, err := s.repo.CreateUser(userData)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *userService) UpdateUser(requestBody map[string]interface{}, userID int) (*response.UserResponse, error) {
	// Use getStringValue to safely handle nil values and type conversions
	var territory_type string
	var territory_id int
	if requestBody["role_id"] != "1" && requestBody["role_id"] != "2" && requestBody["role_id"] != "12" {
		// Role Area
		if requestBody["role_id"] == "3" {
			territory_type = "App\\Models\\Area"
			if v, ok := requestBody["area_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "4" {
			territory_type = "App\\Models\\Region"
			if v, ok := requestBody["region_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "5" || requestBody["role_id"] == "7" || requestBody["role_id"] == "9" || requestBody["role_id"] == "10" || requestBody["role_id"] == "11" {
			territory_type = "App\\Models\\Branch"
			if v, ok := requestBody["branch_id"].(float64); ok {
				territory_id = int(v)
			}
		} else if requestBody["role_id"] == "6" || requestBody["role_id"] == "8" {
			territory_type = "App\\Models\\Cluster"
			if v, ok := requestBody["cluster_id"].(float64); ok {
				territory_id = int(v)
			}
		}
	}
	userData := map[string]string{
		"name":              getStringValue(requestBody["name"]),
		"email":             getStringValue(requestBody["email"]),
		"msisdn":            getStringValue(requestBody["msisdn"]),
		"user_status":       "active",
		"user_type":         getStringValue(requestBody["user_type"]),
		"territory_type":    territory_type,
		"territory_id":      getStringValue(territory_id),
		"outlet_id_digipos": getStringValue(requestBody["outlet_id_digipos"]),
		"nami_agent_id":     getStringValue(requestBody["nami_agent_id"]),
		"password":          getStringValue(requestBody["password"]),
	}

	accounts, err := s.repo.UpdateUser(userData, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
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

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
