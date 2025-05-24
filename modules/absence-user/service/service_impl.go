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
	absences, total, err := s.repo.GetAllAbsences(limit, paginate, page, filters, user_id, month, year, absence_type)
	if err != nil {
		return nil, 0, err
	}

	baseURL := os.Getenv("BASE_URL") // Ambil dari environment variable

	for i := range absences {
		if absences[i].EvidenceImage != nil && *absences[i].EvidenceImage != "" {
			// Ganti \ dengan /
			if absences[i].EvidenceImage != nil {
				updatedValue := strings.ReplaceAll(*absences[i].EvidenceImage, "\\", "/")
				absences[i].EvidenceImage = &updatedValue
			}

			// Tambahkan BASE_URL jika belum ada http/https
			if !strings.HasPrefix(*absences[i].EvidenceImage, "http") {
				updatedValue := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(*absences[i].EvidenceImage, "/"))
				absences[i].EvidenceImage = &updatedValue
			}
		}
	}

	return absences, total, nil
}

func (s *absenceUserService) GetAbsenceUserByID(id int) (*response.ResponseSingleAbsenceUser, error) {
	absenceUser, err := s.repo.GetAbsenceUserByID(id)
	if err != nil {
		return nil, err
	}

	var detailVisitMap *map[string]string
	var targetMap *response.OrderedTargetMap = nil
	baseURL := os.Getenv("BASE_URL")

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

	if absenceUser.EvidenceImage != nil && *absenceUser.EvidenceImage != "" {
		val := strings.ReplaceAll(*absenceUser.EvidenceImage, "\\", "/")

		// Tambahkan BASE_URL jika belum ada prefix http
		if !strings.HasPrefix(val, "http") {
			val = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(val, "/"))
		}
		absenceUser.EvidenceImage = &val
	}

	res := &response.ResponseSingleAbsenceUser{
		ID:            absenceUser.ID,
		UserID:        absenceUser.UserID,
		SubjectType:   absenceUser.SubjectType,
		SubjectID:     absenceUser.SubjectID,
		Type:          absenceUser.Type,
		ClockIn:       absenceUser.ClockIn,
		ClockOut:      absenceUser.ClockOut,
		Description:   absenceUser.Description,
		Longitude:     absenceUser.Longitude,
		Latitude:      absenceUser.Latitude,
		EvidenceImage: absenceUser.EvidenceImage,
		Status:        absenceUser.Status,
		CreatedAt:     absenceUser.CreatedAt,
		UpdatedAt:     absenceUser.UpdatedAt,
		Account:       absenceUser.Account,
		VisitHistory:  absenceUser.VisitHistory,
		Target:        targetMap,
		DetailVisit:   detailVisitMap,
		UserName:      &absenceUser.UserName,
	}

	return res, nil
}

func (s *absenceUserService) GetAbsenceUserToday(only_today bool, user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, string, error) {
	return s.repo.GetAbsenceUserToday(only_today, user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}

func (s *absenceUserService) CreateAbsenceUser(user_id int, subject_type string, subject_id int, description *string, type_absence *string, latitude *string, longitude *string, status *uint, evidenceImage *string) (*models.AbsenceUser, error) {
	convertedUserID := uint(user_id)
	AbsenceUser := &models.AbsenceUser{
		UserID:        &convertedUserID,
		SubjectType:   &subject_type,
		SubjectID:     func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description:   *description,
		Type:          type_absence,
		Latitude:      *latitude,
		Longitude:     *longitude,
		ClockIn:       time.Now(),
		ClockOut:      nil, // ClockOut is now a pointer to time.Time
		Status:        status,
		EvidenceImage: evidenceImage,
	}
	return s.repo.CreateAbsenceUser(AbsenceUser)
}

func (s *absenceUserService) UpdateAbsenceUser(absence_id int, user_id int, subject_type string, subject_id int, description *string, type_absence *string, status *uint) (*models.AbsenceUser, error) {
	AbsenceUser := &models.AbsenceUser{
		ID:          uint(absence_id),
		UserID:      func(v int) *uint { u := uint(v); return &u }(user_id),
		SubjectType: &subject_type,
		SubjectID:   func(v int) *uint { u := uint(v); return &u }(subject_id),
		Description: *description,
		Type:        type_absence,
		ClockOut:    func(t time.Time) *time.Time { return &t }(time.Now()),
		Status:      status,
	}
	return s.repo.UpdateAbsenceUser(AbsenceUser, absence_id)
}

func (s *absenceUserService) GetAbsenceActive(user_id int, type_absence string) ([]models.AbsenceUser, error) {
	return s.repo.GetAbsenceActive(user_id, type_absence)
}

func (s *absenceUserService) AlreadyAbsenceInSameDay(user_id int, type_absence *string, type_checking string, action_type string, subject_type string, subject_id int) (*models.AbsenceUser, error) {
	return s.repo.AlreadyAbsenceInSameDay(user_id, type_absence, type_checking, action_type, subject_type, subject_id)
}

func (s *absenceUserService) DeleteAbsenceUser(id int) error {
	return s.repo.DeleteAbsenceUser(id)
}
