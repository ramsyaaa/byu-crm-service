package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-schedule/repository"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type accountScheduleService struct {
	repo repository.AccountScheduleRepository
}

func NewAccountScheduleService(repo repository.AccountScheduleRepository) AccountScheduleService {
	return &accountScheduleService{repo: repo}
}

func (s *accountScheduleService) GetBySubject(subject_type string, account_id uint) ([]models.AccountSchedule, error) {
	return s.repo.GetBySubject(subject_type, account_id)
}

func (s *accountScheduleService) Insert(requestBody map[string]interface{}, subject_type string, subject_id uint) ([]models.AccountSchedule, error) {
	// Delete existing Schedule for the given subject_type and subject_id
	if err := s.repo.DeleteBySubject(subject_type, subject_id); err != nil {
		return nil, err
	}

	dateSchedule, exists := requestBody["date"]
	if !exists {
		return nil, errors.New("date is missing")
	}

	var DataDateSchedule []time.Time

	switch v := dateSchedule.(type) {
	case string: // Jika hanya satu nilai string, langsung parsing
		v = strings.TrimSpace(v)
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		DataDateSchedule = append(DataDateSchedule, parsedTime)

	case []string: // Jika sudah array string, langsung parsing setiap elemen
		for _, dt := range v {
			dt = strings.TrimSpace(dt)
			parsedTime, err := time.Parse("2006-01-02", dt)
			if err != nil {
				return nil, fmt.Errorf("error parsing date: %v", err)
			}
			DataDateSchedule = append(DataDateSchedule, parsedTime)
		}

	case []interface{}: // Jika array bertipe []interface{}, konversi ke string lalu parsing
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("date contains non-string value")
			}
			strVal = strings.TrimSpace(strVal)
			parsedTime, err := time.Parse("2006-01-02", strVal)
			if err != nil {
				return nil, fmt.Errorf("error parsing date: %v", err)
			}
			DataDateSchedule = append(DataDateSchedule, parsedTime)
		}

	default:
		return nil, fmt.Errorf("invalid date type: %v", reflect.TypeOf(dateSchedule))
	}

	title, exists := requestBody["title"]
	if !exists {
		return nil, errors.New("title is missing")
	}

	var DataTitle []string

	switch v := title.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataTitle = append(DataTitle, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataTitle
		DataTitle = append(DataTitle, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("title contains non-string value")
			}
			DataTitle = append(DataTitle, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid title type: %v", reflect.TypeOf(title))
	}

	schedule_category, exists := requestBody["schedule_category"]
	if !exists {
		return nil, errors.New("schedule category is missing")
	}

	var DataScheduleCategory []string

	switch v := schedule_category.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataScheduleCategory = append(DataScheduleCategory, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataScheduleCategory
		DataScheduleCategory = append(DataScheduleCategory, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("schedule category contains non-string value")
			}
			DataScheduleCategory = append(DataScheduleCategory, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid schedule category type: %v", reflect.TypeOf(schedule_category))
	}

	if len(DataTitle) != len(DataScheduleCategory) && len(DataTitle) != len(DataDateSchedule) {
		return nil, errors.New("title and schedule category length mismatch")
	}

	var insertAccountSchedule []models.AccountSchedule
	for i := range DataTitle {
		insertAccountSchedule = append(insertAccountSchedule, models.AccountSchedule{
			SubjectType:      &subject_type,
			SubjectID:        &subject_id,
			Date:             &DataDateSchedule[i],
			Title:            &DataTitle[i],
			ScheduleCategory: &DataScheduleCategory[i],
		})
	}

	// Insert the new Schedule accounts into the database
	if err := s.repo.Insert(insertAccountSchedule); err != nil {
		return nil, err
	}

	return insertAccountSchedule, nil
}
