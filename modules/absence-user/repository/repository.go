package repository

import "byu-crm-service/models"

type AbsenceUserRepository interface {
	GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error)
	GetAbsenceUserByID(id int) (*models.AbsenceUser, error)
	GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error)
	CreateAbsenceUser(AbsenceUser *models.AbsenceUser) (*models.AbsenceUser, error)
	UpdateAbsenceUser(AbsenceUser *models.AbsenceUser, id int) (*models.AbsenceUser, error)
}
