package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type visitHistoryRepository struct {
	db *gorm.DB
}

func NewVisitHistoryRepository(db *gorm.DB) VisitHistoryRepository {
	return &visitHistoryRepository{db: db}
}

func (r *visitHistoryRepository) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error) {
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

func (r *visitHistoryRepository) GetAbsenceUserByID(id int) (*models.AbsenceUser, error) {
	var absence_user models.AbsenceUser
	err := r.db.First(&absence_user, id).Error
	if err != nil {
		return nil, err
	}
	return &absence_user, nil
}

func (r *visitHistoryRepository) GetAbsenceUserToday(
	only_today bool,
	user_id int,
	type_absence *string,
	type_checking string,
	action_type string,
	subject_type string,
	subject_id int,
) (*models.AbsenceUser, string, error) {
	var absence_user models.AbsenceUser

	query := r.db.Where("user_id = ?", user_id)

	// Default message
	message := "Absence user found"

	// Filter by type_checking
	if only_today {
		switch type_checking {
		case "daily":
			query = query.Where("DATE(clock_in) = CURDATE()")
			message = "The absence user is already checked in today"
		case "monthly":
			query = query.Where("MONTH(clock_in) = MONTH(CURDATE()) AND YEAR(clock_in) = YEAR(CURDATE())")
			message = "The absence user is already checked in this month"
		}
		switch action_type {
		case "Clock Out":
			query = query.Where("clock_out IS NOT NULL")
		}
	} else {
		// Filter by action_type
		switch action_type {
		case "Clock In":
			query = query.Where("clock_in IS NOT NULL AND clock_out IS NULL")
			message = "The user already clocked in, please clock out first"
		case "Clock Out":
			query = query.Where("clock_out IS NULL")
			message = "The user already clocked out, need to clock in first"
		}
	}

	// Filter by type_absence
	if type_absence != nil {
		query = query.Where("type = ?", *type_absence)
	}

	// Filter by subject_type and subject_id if provided
	if subject_type != "" && subject_id != 0 {

		query = query.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id)
	}

	err := query.First(&absence_user).Error
	if err != nil {
		return nil, message, err
	}

	return &absence_user, message, nil
}

func (r *visitHistoryRepository) CreateVisitHistory(visit_history *models.VisitHistory) (*models.VisitHistory, error) {
	if err := r.db.Create(visit_history).Error; err != nil {
		return nil, err
	}

	var createdVisitHistory models.VisitHistory
	if err := r.db.First(&createdVisitHistory, "id = ?", visit_history.ID).Error; err != nil {
		return nil, err
	}

	return &createdVisitHistory, nil
}

func (r *visitHistoryRepository) UpdateAbsenceUser(absence_user *models.AbsenceUser, id int) (*models.AbsenceUser, error) {
	var existingAbsenceUser models.AbsenceUser
	if err := r.db.First(&existingAbsenceUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingAbsenceUser).Updates(absence_user).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingAbsenceUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &existingAbsenceUser, nil
}
