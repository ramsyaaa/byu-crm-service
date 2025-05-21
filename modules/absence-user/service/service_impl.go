package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/absence-user/repository"
	"byu-crm-service/modules/absence-user/response"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type absenceUserService struct {
	repo repository.AbsenceUserRepository
}

func NewAbsenceUserService(repo repository.AbsenceUserRepository) AbsenceUserService {
	return &absenceUserService{repo: repo}
}

func (s *absenceUserService) GetAllAbsences(limit int, paginate bool, page int, filters map[string]string, user_id int, month int, year int, absence_type string) ([]models.AbsenceUser, int64, error) {
	return s.repo.GetAllAbsences(limit, paginate, page, filters, user_id, month, year, absence_type)
}

func (s *absenceUserService) GetAbsenceUserByID(id int) (*response.ResponseSingleAbsenceUser, error) {
	absenceUser, err := s.repo.GetAbsenceUserByID(id)
	if err != nil {
		return nil, err
	}

	var detailVisitMap *map[string]string
	var targetMap *response.OrderedTargetMap = nil

	if absenceUser.VisitHistory != nil {
		// Parse target
		if absenceUser.VisitHistory.Target != nil {
			var tempTarget map[string]int
			if err := json.Unmarshal([]byte(*absenceUser.VisitHistory.Target), &tempTarget); err == nil {
				targetMap = &response.OrderedTargetMap{Data: tempTarget}
			}
		}

		// Parse detail_visit
		if absenceUser.VisitHistory.DetailVisit != nil {
			var tempDetailVisit map[string]string
			if err := json.Unmarshal([]byte(*absenceUser.VisitHistory.DetailVisit), &tempDetailVisit); err == nil {
				baseURL := os.Getenv("BASE_URL")

				// Tambahkan prefix base URL jika ada file presentasi_demo atau dealing_sekolah
				for key, val := range tempDetailVisit {
					if key == "presentasi_demo" || key == "dealing_sekolah" {
						if val != "" {
							// Ganti semua backslash jadi slash agar menjadi path yang valid
							val = strings.ReplaceAll(val, "\\", "/")

							// Tambahkan BASE_URL jika belum ada prefix http
							if !strings.HasPrefix(val, "http") {
								val = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(val, "/"))
							}
							tempDetailVisit[key] = val
						}
					}
				}

				detailVisitMap = &tempDetailVisit
			}
		}

	}

	res := &response.ResponseSingleAbsenceUser{
		ID:           absenceUser.ID,
		UserID:       absenceUser.UserID,
		SubjectType:  absenceUser.SubjectType,
		SubjectID:    absenceUser.SubjectID,
		Type:         absenceUser.Type,
		ClockIn:      absenceUser.ClockIn,
		ClockOut:     absenceUser.ClockOut,
		Description:  absenceUser.Description,
		Longitude:    absenceUser.Longitude,
		Latitude:     absenceUser.Latitude,
		CreatedAt:    absenceUser.CreatedAt,
		UpdatedAt:    absenceUser.UpdatedAt,
		Account:      absenceUser.Account,
		VisitHistory: absenceUser.VisitHistory,
		Target:       targetMap,
		DetailVisit:  detailVisitMap,
	}

	return res, nil
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

func (s *absenceUserService) UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string) (*models.AbsenceUser, error) {
	AbsenceUser := &models.AbsenceUser{
		ID:          uint(absence_id),
		UserID:      func(v int) *uint { u := uint(v); return &u }(user_id),
		SubjectType: &subject_type,
		SubjectID:   func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description: *description,
		Type:        type_absence,
		ClockOut:    func(t time.Time) *time.Time { return &t }(time.Now()),
	}
	return s.repo.UpdateAbsenceUser(AbsenceUser, absence_id)
}

func (s *absenceUserService) GetAbsenceActive(user_id int, type_absence string) ([]models.AbsenceUser, error) {
	return s.repo.GetAbsenceActive(user_id, type_absence)
}

func (s *absenceUserService) AlreadyAbsenceInSameDay(user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, error) {
	return s.repo.AlreadyAbsenceInSameDay(user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}
