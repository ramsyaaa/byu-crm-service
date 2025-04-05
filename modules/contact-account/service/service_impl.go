package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/repository"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type contactAccountService struct {
	repo repository.ContactAccountRepository
}

func NewContactAccountService(repo repository.ContactAccountRepository) ContactAccountService {
	return &contactAccountService{repo: repo}
}

func (s *contactAccountService) GetContactAccountByAccountID(account_id uint) ([]models.ContactAccount, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *contactAccountService) InsertContactAccount(requestBody map[string]interface{}, account_id uint) ([]models.ContactAccount, error) {
	// Delete existing contact accounts for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	contactID, exists := requestBody["contact_id"]
	if !exists {
		return nil, errors.New("contact_id is missing")
	}

	var DataContactID []string

	switch v := contactID.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataContactID = append(DataContactID, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataContactID
		DataContactID = append(DataContactID, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("contact_id contains non-string value")
			}
			DataContactID = append(DataContactID, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid contact_id type: %v", reflect.TypeOf(contactID))
	}

	var contactAccounts []models.ContactAccount

	// Loop through the contact IDs and create ContactAccount instances
	for _, contact := range DataContactID {
		idUint, err := strconv.ParseUint(contact, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting contact ID to uint: %v", err)
		}

		contactAccounts = append(contactAccounts, models.ContactAccount{
			ContactID: func(u uint) *uint { return &u }(uint(idUint)),
			AccountID: &account_id,
		})
	}

	// Insert the new contact accounts into the database
	if err := s.repo.Insert(contactAccounts); err != nil {
		return nil, err
	}

	return contactAccounts, nil
}
