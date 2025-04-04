package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type absenceUserRepository struct {
	db *gorm.DB
}

func NewAbsenceUserRepository(db *gorm.DB) AbsenceUserRepository {
	return &absenceUserRepository{db: db}
}

func (r *absenceUserRepository) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error) {
	var absence_users []models.AbsenceUser
	var total int64

	query := r.db.Model(&models.AbsenceUser{})

	if user_id != 0 {
		query = query.Where("absence_users.user_id = ?", user_id)
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("absence_users.description LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("absence_users.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("absence_users.created_at <= ?", endDate)
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

	err := query.Find(&absence_users).Error
	return absence_users, total, err
}

func (r *absenceUserRepository) GetAbsenceUserByID(id int) (*models.AbsenceUser, error) {
	var absence_user models.AbsenceUser
	err := r.db.First(&absence_user, id).Error
	if err != nil {
		return nil, err
	}
	return &absence_user, nil
}

func (r *absenceUserRepository) GetAbsenceUserToday(user_id int, type_absence *string, type_checking string) (*models.AbsenceUser, string, error) {
	var absence_user models.AbsenceUser

	query := r.db.Where("user_id = ?", user_id)

	// Filter berdasarkan type_checking
	var message string
	if type_checking == "daily" {
		query = query.Where("DATE(date) = CURDATE()")
		message = "The absence user is already checked in today"
	} else if type_checking == "monthly" {
		message = "The absence user is already checked in this month"
		query = query.Where("MONTH(date) = MONTH(CURDATE()) AND YEAR(date) = YEAR(CURDATE())")
	}

	// Filter berdasarkan type_absence jika ada
	if type_absence != nil {
		query = query.Where("type = ?", *type_absence)
	}

	err := query.First(&absence_user).Error
	if err != nil {
		return nil, message, err
	}

	return &absence_user, message, nil
}

func (r *absenceUserRepository) CreateAbsenceUser(absence_user *models.AbsenceUser) (*models.AbsenceUser, error) {
	if err := r.db.Create(absence_user).Error; err != nil {
		return nil, err
	}

	var createdAbsenceUser models.AbsenceUser
	if err := r.db.First(&createdAbsenceUser, "id = ?", absence_user.ID).Error; err != nil {
		return nil, err
	}

	return &createdAbsenceUser, nil
}

func (r *absenceUserRepository) UpdateAbsenceUser(faculty *models.AbsenceUser, id int) (*models.AbsenceUser, error) {
	var existingFaculty models.AbsenceUser
	if err := r.db.First(&existingFaculty, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingFaculty).Updates(faculty).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingFaculty, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &existingFaculty, nil
}
