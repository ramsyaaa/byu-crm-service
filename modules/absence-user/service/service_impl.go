package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/absence-user/repository"
	"time"
)

type absenceUserService struct {
	repo repository.AbsenceUserRepository
}

func NewAbsenceUserService(repo repository.AbsenceUserRepository) AbsenceUserService {
	return &absenceUserService{repo: repo}
}

func (s *absenceUserService) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error) {
	return s.repo.GetAllAbsences(limit, paginate, page, filters, user_id)
}

func (s *absenceUserService) GetAbsenceUserByID(id int) (*models.AbsenceUser, error) {
	return s.repo.GetAbsenceUserByID(id)
}

func (s *absenceUserService) GetAbsenceUserToday(user_id int, type_absence *string, type_checking string) (*models.AbsenceUser, string, error) {
	return s.repo.GetAbsenceUserToday(user_id, type_absence, type_checking)
}

func (s *absenceUserService) CreateAbsenceUser(user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error) {
	convertedUserID := uint(user_id)
	AbsenceUser := &models.AbsenceUser{
		UserID:      &convertedUserID,
		SubjectType: &subject_type,
		SubjectID:   func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description: *description,
		Type:        type_absence,
		Latitude:    *latitude,
		Longitude:   *longitude,
		Date:        time.Now(),
	}
	return s.repo.CreateAbsenceUser(AbsenceUser)
}

func (s *absenceUserService) UpdateAbsenceUser(name *string, id int) (*models.AbsenceUser, error) {
	AbsenceUser := &models.AbsenceUser{}
	return s.repo.UpdateAbsenceUser(AbsenceUser, id)
}
