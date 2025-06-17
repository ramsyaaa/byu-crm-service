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

func (r *absenceUserRepository) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int, month int, year int, absence_type string) ([]models.AbsenceUser, int64, error) {
	var absence_users []models.AbsenceUser
	var total int64

	query := r.db.Model(&models.AbsenceUser{}).
		Joins("LEFT JOIN users ON users.id = absence_users.user_id").
		Joins("LEFT JOIN accounts ON accounts.id = absence_users.subject_id AND absence_users.subject_type = ?", "App\\Models\\Account").
		Joins("LEFT JOIN cities ON cities.id = accounts.city").
		Joins("LEFT JOIN account_type_school_details ON account_type_school_details.account_id = accounts.id").
		Select(`absence_users.*, 
			users.name AS user_name, 
			users.yae_code AS yae_code,
			cities.name AS city_name,
			cities.cluster_id AS cluster_id,
			cities.id AS city_id,
			accounts.account_name, 
			accounts.account_code, 
			account_type_school_details.dies_natalis, 
			account_type_school_details.extracurricular, 
			account_type_school_details.football_field_brannnding, 
			account_type_school_details.basketball_field_branding, 
			account_type_school_details.wall_painting_branding, 
			account_type_school_details.wall_magazine_branding`)

	// Filter: all_user
	if allUser, exists := filters["all_user"]; exists && allUser != "1" {
		if user_id != 0 {
			query = query.Where("absence_users.user_id = ?", user_id)
		}
	}

	// Filter: search
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("absence_users.description LIKE ? OR users.name LIKE ?", "%"+token+"%", "%"+token+"%"),
			)
		}
	}

	// Filter: date range
	startDate, hasStart := filters["start_date"]
	endDate, hasEnd := filters["end_date"]

	if hasStart && startDate != "" && hasEnd && endDate != "" {
		startDateTime := startDate + " 00:00:00"
		endDateTime := endDate + " 23:59:59"
		query = query.Where("absence_users.created_at BETWEEN ? AND ?", startDateTime, endDateTime)
	} else if hasStart && startDate != "" {
		startDateTime := startDate + " 00:00:00"
		query = query.Where("absence_users.created_at >= ?", startDateTime)
	} else if hasEnd && endDate != "" {
		endDateTime := endDate + " 23:59:59"
		query = query.Where("absence_users.created_at <= ?", endDateTime)
	}

	// Filter: month & year
	if month != 0 && year != 0 {
		query = query.Where("MONTH(absence_users.created_at) = ? AND YEAR(absence_users.created_at) = ?", month, year)
	}

	// Filter: absence_type
	if absence_type != "" {
		query = query.Where("absence_users.type = ?", absence_type)
	}

	// Filter: status
	if status, hasStatus := filters["status"]; hasStatus && status == "0" {
		query = query.Where("absence_users.status = ?", status)
	}

	// Get total count before pagination
	query.Count(&total)

	// Ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy == "" {
		orderBy = "absence_users.created_at"
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&absence_users).Error
	if err != nil {
		return nil, total, err
	}

	// Optional: Load full Account with SchoolDetail for each absence_user if needed
	for i := range absence_users {
		if absence_users[i].SubjectType != nil && *absence_users[i].SubjectType == "App\\Models\\Account" {
			var account models.Account
			if err := r.db.First(&account, absence_users[i].SubjectID).Error; err == nil {
				absence_users[i].Account = &account
			}
		}

		// Join visit_histories
		var visit models.VisitHistory
		if err := r.db.Where("absence_user_id = ?", absence_users[i].ID).First(&visit).Error; err == nil {
			absence_users[i].VisitHistory = &visit
		}
	}

	return absence_users, total, nil
}

func (r *absenceUserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.AbsenceUser{}).Where("id = ?", id).Updates(fields).Error
}

func (r *absenceUserRepository) GetAbsenceUserByID(id int) (*models.AbsenceUser, error) {
	var absence_user models.AbsenceUser

	err := r.db.
		Joins("LEFT JOIN users ON users.id = absence_users.user_id").
		Select("absence_users.*, users.name AS user_name").
		First(&absence_user, id).Error
	if err != nil {
		return nil, err
	}

	// Load related Account jika SubjectType = Account
	if absence_user.SubjectType != nil && absence_user.SubjectID != nil && *absence_user.SubjectType == "App\\Models\\Account" {
		var account models.Account
		if err := r.db.First(&account, *absence_user.SubjectID).Error; err == nil {
			absence_user.Account = &account
		}
	}

	// Load related VisitHistory
	var visit models.VisitHistory
	if err := r.db.Where("absence_user_id = ?", absence_user.ID).First(&visit).Error; err == nil {
		absence_user.VisitHistory = &visit
	}

	return &absence_user, nil
}

func (r *absenceUserRepository) GetAbsenceUserToday(
	only_today bool,
	user_id int,
	type_absence *string,
	type_checking string,
	action_type string,
	subject_type string,
	subject_id int,
) (*models.AbsenceUser, string, error) {
	var absence_user models.AbsenceUser

	// Always declare query and assign base query
	query := r.db.Model(&models.AbsenceUser{})
	if user_id != 0 {
		query = query.Where("user_id = ?", user_id)
	}

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
	}

	switch action_type {
	case "Clock In":
		query = query.Where("clock_in IS NOT NULL AND clock_out IS NULL")
		message = "The user already clocked in, please clock out first"
		if user_id == 0 {
			message = "Terdapat user lain yang sudah clock in"
		}
	case "Clock Out":
		query = query.Where("clock_out IS NULL")
		message = "The user already clocked out, need to clock in first"
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

func (r *absenceUserRepository) UpdateAbsenceUser(absence_user *models.AbsenceUser, id int) (*models.AbsenceUser, error) {
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

func (r *absenceUserRepository) GetAbsenceActive(user_id int, type_absence string) ([]models.AbsenceUser, error) {
	var absence_users []models.AbsenceUser

	query := r.db.Where("user_id = ? AND (clock_out IS NULL OR DATE(clock_in) = CURRENT_DATE)", user_id)
	if type_absence != "" {
		query = query.Where("type = ?", type_absence)
	}

	err := query.Find(&absence_users).Error
	if err != nil {
		return nil, err
	}

	// Ambil data subject berdasarkan subject_type dan subject_id
	for i, absence := range absence_users {
		switch *absence.SubjectType {
		case "App\\Models\\Account":
			var account models.Account
			if err := r.db.First(&account, absence.SubjectID).Error; err == nil {
				absence_users[i].Account = &account // Ensure the Account field is populated
			}
		}
	}

	return absence_users, nil
}

func (r *absenceUserRepository) AlreadyAbsenceInSameDay(user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, error) {
	var absence_user models.AbsenceUser

	query := r.db.Where("user_id = ?", user_id)

	switch type_checking {
	case "daily":
		query = query.Where("DATE(clock_in) = CURDATE()")
	case "monthly":
		query = query.Where("MONTH(clock_in) = MONTH(CURDATE()) AND YEAR(clock_in) = YEAR(CURDATE())")
	}

	if type_absence != nil {
		query = query.Where("type = ?", *type_absence)
	}

	// Filter by subject_type and subject_id if provided
	if subject_type != "" && subject_id != 0 {

		query = query.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id)
	}

	err := query.First(&absence_user).Error
	if err != nil {
		return nil, err
	}

	return &absence_user, nil
}

func (r *absenceUserRepository) DeleteAbsenceUser(id int) error {
	return r.db.Delete(&models.AbsenceUser{}, id).Error
}
