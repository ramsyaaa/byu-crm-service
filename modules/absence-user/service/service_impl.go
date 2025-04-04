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

func (s *absenceUserService) GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error) {
	return s.repo.GetAbsenceUserToday(only_today, user_id, type_absence, type_checking, action_type, subject_type, subject_id)
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
		ClockIn:     time.Now(),
		ClockOut:    nil, // ClockOut is now a pointer to time.Time
	}
	return s.repo.CreateAbsenceUser(AbsenceUser)
}

func (s *absenceUserService) UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error) {
	AbsenceUser := &models.AbsenceUser{
		ID:          uint(absence_id),
		UserID:      func(v int) *uint { u := uint(v); return &u }(user_id),
		SubjectType: &subject_type,
		SubjectID:   func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description: *description,
		Type:        type_absence,
		Latitude:    *latitude,
		Longitude:   *longitude,
		ClockOut:    func(t time.Time) *time.Time { return &t }(time.Now()),
	}
	return s.repo.UpdateAbsenceUser(AbsenceUser, absence_id)
}
