package repository

import (
	"byu-crm-service/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Msisdn        string `json:"msisdn"`
	UserStatus    string `json:"user_status"`
	UserType      string `json:"user_type"`
	TerritoryID   uint   `json:"territory_id"`
	TerritoryType string `json:"territory_type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(limit int, paginate bool, page int, filters map[string]string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("users.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("users.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("users.created_at <= ?", endDate)
	}

	// Get total count before applying pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Apply pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByID(id uint) (*UserResponse, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Map hanya field yang diperlukan
	response := &UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Msisdn:        user.Msisdn,
		UserStatus:    user.UserStatus,
		UserType:      user.UserType,
		TerritoryID:   user.TerritoryID,
		TerritoryType: user.TerritoryType,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	return response, nil
}

func (r *userRepository) UpdateUserProfile(id uint, user map[string]interface{}) (*UserResponse, error) {
	var existingUser models.User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	existingUser.Password = user["password"].(string)

	if err := r.db.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	response := &UserResponse{
		ID:            existingUser.ID,
		Name:          existingUser.Name,
		Email:         existingUser.Email,
		Avatar:        existingUser.Avatar,
		Msisdn:        existingUser.Msisdn,
		UserStatus:    existingUser.UserStatus,
		UserType:      existingUser.UserType,
		TerritoryID:   existingUser.TerritoryID,
		TerritoryType: existingUser.TerritoryType,
		CreatedAt:     existingUser.CreatedAt,
		UpdatedAt:     existingUser.UpdatedAt,
	}

	return response, nil
}
