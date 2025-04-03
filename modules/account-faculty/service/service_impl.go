package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-faculty/repository"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type accountFacultyService struct {
	repo repository.AccountFacultyRepository
}

func NewAccountFacultyService(repo repository.AccountFacultyRepository) AccountFacultyService {
	return &accountFacultyService{repo: repo}
}

func (s *accountFacultyService) GetByAccountID(account_id uint) ([]models.AccountFaculty, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *accountFacultyService) Insert(requestBody map[string]interface{}, account_id uint) ([]models.AccountFaculty, error) {
	// Delete existing faculty for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	faculties, exists := requestBody["faculties"]
	if !exists {
		return nil, errors.New("faculties is missing")
	}

	var DataFaculty []string

	switch v := faculties.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataFaculty = append(DataFaculty, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataFaculty
		DataFaculty = append(DataFaculty, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("contact_id contains non-string value")
			}
			DataFaculty = append(DataFaculty, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid contact_id type: %v", reflect.TypeOf(faculties))
	}

	var insertFaculties []models.AccountFaculty
	for i := range DataFaculty {
		// Conversion to uint64
		// Assuming DataFaculty[i] is a string that can be converted to uint64
		facultyID, err := strconv.ParseUint(DataFaculty[i], 10, 64)
		if err != nil {
			fmt.Println("Error converting FacultyID:", err)
			continue // Skipped facultyID if conversion fails
		}

		// Conversion to *uint
		facultyIDUint := uint(facultyID)

		insertFaculties = append(insertFaculties, models.AccountFaculty{
			AccountID: &account_id,
			FacultyID: &facultyIDUint,
		})
	}

	// Insert the new contact accounts into the database
	if err := s.repo.Insert(insertFaculties); err != nil {
		return nil, err
	}

	return insertFaculties, nil
}
