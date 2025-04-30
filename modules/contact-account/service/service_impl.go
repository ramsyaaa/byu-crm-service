package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/repository"
	"byu-crm-service/modules/contact-account/response"
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

func (s *contactAccountService) GetAllContacts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, accountID int) ([]response.ContactResponse, int64, error) {
	return s.repo.GetAllContacts(limit, paginate, page, filters, userRole, territoryID, accountID)
}

func (s *contactAccountService) GetContactAccountByAccountID(account_id uint) ([]models.ContactAccount, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *contactAccountService) FindByContactID(id uint, userRole string, territoryID uint) (*response.SingleContactResponse, error) {
	contact, err := s.repo.FindByContactID(id, userRole, territoryID)
	if err != nil {
		return nil, err
	}
	var contactResponse response.SingleContactResponse
	contactResponse.ID = contact.ID
	contactResponse.ContactName = contact.ContactName
	contactResponse.PhoneNumber = contact.PhoneNumber
	contactResponse.Position = contact.Position
	contactResponse.Birthday = contact.Birthday
	contactResponse.CreatedAt = contact.CreatedAt
	contactResponse.UpdatedAt = contact.UpdatedAt
	contactResponse.SocialMedias = contact.SocialMedias
	contactResponse.Accounts = contact.Accounts

	contactResponse.Category = []string{}
	contactResponse.Url = []string{}

	for _, sm := range contactResponse.SocialMedias {
		contactResponse.Category = append(contactResponse.Category, *sm.Category)
		contactResponse.Url = append(contactResponse.Url, *sm.Url)
	}

	return &contactResponse, nil
}

func (s *contactAccountService) CreateContact(requestBody map[string]interface{}) (*models.Contact, error) {
	contactData := map[string]string{
		"contact_name": requestBody["contact_name"].(string),
		"phone_number": requestBody["phone_number"].(string),
		"position":     requestBody["position"].(string),
		"birthday":     requestBody["birthday"].(string),
	}

	contact, err := s.repo.CreateContact(contactData)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *contactAccountService) UpdateContact(requestBody map[string]interface{}, contactID int) (*models.Contact, error) {
	contactData := map[string]string{
		"contact_name": requestBody["contact_name"].(string),
		"phone_number": requestBody["phone_number"].(string),
		"position":     requestBody["position"].(string),
		"birthday":     requestBody["birthday"].(string),
	}

	contact, err := s.repo.UpdateContact(contactData, contactID)
	if err != nil {
		return nil, err
	}

	return contact, nil
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

func (s *contactAccountService) InsertContactAccountByContactID(requestBody map[string]interface{}, contact_id uint) ([]models.ContactAccount, error) {
	// Delete existing contact accounts for the given contact_id
	if err := s.repo.DeleteAccountByContactID(contact_id); err != nil {
		return nil, err
	}

	accountID, exists := requestBody["account_id"]
	if !exists {
		return nil, errors.New("account_id is missing")
	}

	var DataAccountID []string

	switch v := accountID.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataAccountID = append(DataAccountID, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataAccountID
		DataAccountID = append(DataAccountID, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("account_id contains non-string value")
			}
			DataAccountID = append(DataAccountID, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid account_id type: %v", reflect.TypeOf(accountID))
	}

	var contactAccounts []models.ContactAccount

	// Loop through the account IDs and create ContactAccount instances
	for _, account := range DataAccountID {
		idUint, err := strconv.ParseUint(account, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting account ID to uint: %v", err)
		}

		contactAccounts = append(contactAccounts, models.ContactAccount{
			AccountID: func(u uint) *uint { return &u }(uint(idUint)),
			ContactID: &contact_id,
		})
	}

	// Insert the new contact accounts into the database
	if err := s.repo.Insert(contactAccounts); err != nil {
		return nil, err
	}

	return contactAccounts, nil
}
