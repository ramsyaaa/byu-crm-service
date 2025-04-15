package service

import "byu-crm-service/models"

type AbsenceUserService interface {
	GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int, month int, year int, absence_type string) ([]models.AbsenceUser, int64, error)
	GetAbsenceUserByID(id int) (*models.AbsenceUser, error)
	GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error)
	AlreadyAbsenceInSameDay(user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, error)
	CreateAbsenceUser(user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error)
	UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error)
	GetAbsenceActive(user_id int, type_absence string) ([]models.AbsenceUser, error)
}
