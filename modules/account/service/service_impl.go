package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/repository"
	cityRepository "byu-crm-service/modules/city/repository"
	"errors"
	"fmt"
)

type accountService struct {
	repo     repository.AccountRepository
	cityRepo cityRepository.CityRepository
}

func NewAccountService(repo repository.AccountRepository, cityRepo cityRepository.CityRepository) AccountService {
	return &accountService{repo: repo, cityRepo: cityRepo}
}

func (s *accountService) GetAllAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, map[string]interface{}, error) {
	if limit <= 0 || page <= 0 {
		return nil, nil, errors.New("limit and page must be greater than 0")
	}

	// Call repository layer
	accounts, totalRecords, err := s.repo.GetFilteredAccounts(limit, page, search, userRole, territoryID)
	if err != nil {
		return nil, nil, err
	}

	// Create pagination metadata
	totalPages := (totalRecords + limit - 1) / limit
	pagination := map[string]interface{}{
		"current_page": page,
		"total_pages":  totalPages,
		"total_items":  totalRecords,
		"limit":        limit,
	}

	return accounts, pagination, nil
}

func (s *accountService) ProcessAccount(data []string) error {
	existingAccount, err := s.repo.FindByAccountCode(data[8]) // AccountCode
	if err != nil {
		return err
	}
	if existingAccount != nil {
		return nil
	}

	// Assuming you have a cityService with a method FindCityByName
	city, err := s.cityRepo.FindByName(data[4]) // City
	if err != nil {
		return err
	}
	if city == nil {
		return fmt.Errorf("city not found")
	}

	account := models.Account{
		AccountName:             &data[5],
		AccountType:             &data[6],
		AccountCategory:         &data[7],
		AccountCode:             &data[8],
		City:                    &city.ID, // Set city ID
		ContactName:             &data[9],
		EmailAccount:            &data[10],
		WebsiteAccount:          &data[12],
		Potensi:                 &data[11],
		SystemInformasiAkademik: &data[13],
		Ownership:               &data[14],
	}
	return s.repo.Create(&account)
}
