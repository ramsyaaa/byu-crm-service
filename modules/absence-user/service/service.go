package service

import "byu-crm-service/models"

type AbsenceUserService interface {
	GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error)
	GetAbsenceUserByID(id int) (*models.AbsenceUser, error)
	GetAbsenceUserToday(user_id int, type_absence *string, type_checking string) (*models.AbsenceUser, string, error)
	CreateAbsenceUser(user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error)
	UpdateAbsenceUser(name *string, id int) (*models.AbsenceUser, error)
}
