package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-member/repository"
	"errors"
	"fmt"
	"reflect"
)

type accountMemberService struct {
	repo repository.AccountMemberRepository
}

func NewAccountMemberService(repo repository.AccountMemberRepository) AccountMemberService {
	return &accountMemberService{repo: repo}
}

func (s *accountMemberService) GetBySubject(subject_type string, account_id uint) ([]models.AccountMember, error) {
	return s.repo.GetBySubject(subject_type, account_id)
}

func (s *accountMemberService) Insert(requestBody map[string]interface{}, subject_type string, subject_id uint, key1 string, key2 string) ([]models.AccountMember, error) {
	// Delete existing Member for the given subject_type and subject_id
	if err := s.repo.DeleteBySubject(subject_type, subject_id); err != nil {
		return nil, err
	}

	year, exists := requestBody[key1]
	if !exists {
		return nil, errors.New("year is missing")
	}

	var DataYear []string

	switch v := year.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataYear = append(DataYear, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataYear
		DataYear = append(DataYear, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("year contains non-string value")
			}
			DataYear = append(DataYear, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid year type: %v", reflect.TypeOf(year))
	}

	amount, exists := requestBody[key2]
	if !exists {
		return nil, errors.New("amount is missing")
	}

	var DataAmount []string

	switch v := amount.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataAmount = append(DataAmount, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataAmount
		DataAmount = append(DataAmount, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("amount contains non-string value")
			}
			DataAmount = append(DataAmount, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid amount type: %v", reflect.TypeOf(amount))
	}

	if len(DataYear) != len(DataAmount) {
		return nil, errors.New("year and amount length mismatch")
	}

	var insertAccountMember []models.AccountMember
	for i := range DataYear {
		insertAccountMember = append(insertAccountMember, models.AccountMember{
			SubjectType: &subject_type,
			SubjectID:   &subject_id,
			Year:        &DataYear[i],
			Amount:      &DataAmount[i],
		})
	}

	// Insert the new member accounts into the database
	if err := s.repo.Insert(insertAccountMember); err != nil {
		return nil, err
	}

	return insertAccountMember, nil
}
