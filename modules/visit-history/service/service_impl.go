package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/visit-history/repository"
	"time"
)

type visitHistoryService struct {
	repo repository.VisitHistoryRepository
}

func NewVisitHistoryService(repo repository.VisitHistoryRepository) VisitHistoryService {
	return &visitHistoryService{repo: repo}
}

func (s *visitHistoryService) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int) ([]models.AbsenceUser, int64, error) {
	return s.repo.GetAllAbsences(limit, paginate, page, filters, user_id)
}

func (s *visitHistoryService) GetAbsenceUserByID(id int) (*models.AbsenceUser, error) {
	return s.repo.GetAbsenceUserByID(id)
}

func (s *visitHistoryService) GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error) {
	return s.repo.GetAbsenceUserToday(only_today, user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}

func (s *visitHistoryService) CreateVisitHistory(user_id int, subject_type string, subject_id int, absence_user_id int, greeting bool, survey bool, presentation bool, description *string) (*models.VisitHistory, error) {
	convertedUserID := uint(user_id)
	VisitHistory := &models.VisitHistory{
		UserID:        &convertedUserID,
		SubjectType:   &subject_type,
		SubjectID:     func(v int) *uint { u := uint(v); return &u }(subject_id),
		AbsenceUserID: func(v int) *uint { u := uint(v); return &u }(absence_user_id),
		Greeting:      greeting,
		Survey:        survey,
		Presentation:  presentation,
		Description:   *description,
	}
	return s.repo.CreateVisitHistory(VisitHistory)
}

func (s *visitHistoryService) UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string) (*models.AbsenceUser, error) {
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
