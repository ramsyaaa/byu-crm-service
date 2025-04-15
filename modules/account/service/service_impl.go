package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/repository"
	cityRepository "byu-crm-service/modules/city/repository"
	"fmt"
	"strconv"
)

type accountService struct {
	repo     repository.AccountRepository
	cityRepo cityRepository.CityRepository
}

func NewAccountService(repo repository.AccountRepository, cityRepo cityRepository.CityRepository) AccountService {
	return &accountService{repo: repo, cityRepo: cityRepo}
}

func (s *accountService) GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool, excludeVisited bool) ([]models.Account, int64, error) {
	return s.repo.GetAllAccounts(limit, paginate, page, filters, userRole, territoryID, userID, onlyUserPic, excludeVisited)
}

func (s *accountService) CreateAccount(requestBody map[string]interface{}, userID int) ([]models.Account, error) {
	accountData := map[string]string{
		"account_name":              requestBody["account_name"].(string),
		"account_image":             requestBody["account_image"].(string),
		"account_type":              requestBody["account_type"].(string),
		"account_category":          requestBody["account_category"].(string),
		"account_code":              requestBody["account_code"].(string),
		"city":                      requestBody["city"].(string),
		"contact_name":              requestBody["contact_name"].(string),
		"email_account":             requestBody["email_account"].(string),
		"website_account":           requestBody["website_account"].(string),
		"system_informasi_akademik": requestBody["system_informasi_akademik"].(string),
		"ownership":                 requestBody["ownership"].(string),
		"pic":                       requestBody["pic"].(string),
		"pic_internal":              requestBody["pic_internal"].(string),
		"latitude":                  requestBody["latitude"].(string),
		"longitude":                 requestBody["longitude"].(string),
	}

	accounts, err := s.repo.CreateAccount(accountData, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountService) UpdateAccount(requestBody map[string]interface{}, accountID int, userRole string, territoryID int, userID int) ([]models.Account, error) {
	existingAccount, err := s.repo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))

	if err != nil {
		return nil, err
	}

	getString := func(val *string) string {
		if val != nil {
			return *val
		}
		return ""
	}

	accountData := map[string]string{
		"account_name":              getString(existingAccount.AccountName),
		"account_image":             getString(existingAccount.AccountImage),
		"account_type":              getString(existingAccount.AccountType),
		"account_category":          getString(existingAccount.AccountCategory),
		"account_code":              getString(existingAccount.AccountCode),
		"city":                      getString(stringPointer(fmt.Sprintf("%d", *existingAccount.City))),
		"contact_name":              getString(existingAccount.ContactName),
		"email_account":             getString(existingAccount.EmailAccount),
		"website_account":           getString(existingAccount.WebsiteAccount),
		"system_informasi_akademik": getString(existingAccount.SystemInformasiAkademik),
		"ownership":                 getString(existingAccount.Ownership),
		"pic":                       getString(existingAccount.Pic),
		"pic_internal":              getString(existingAccount.PicInternal),
		"latitude":                  getString(existingAccount.Latitude),
		"longitude":                 getString(existingAccount.Longitude),
	}

	// Override field hanya jika tersedia di requestBody
	for key := range accountData {
		if val, ok := requestBody[key]; ok && val != nil {
			accountData[key] = fmt.Sprintf("%v", val)
		}
	}

	accounts, err := s.repo.UpdateAccount(accountData, accountID, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountService) GetAccountVisitCounts(filters map[string]string, userRole string, territoryID int, userID int) (int64, int64, int64, error) {
	return s.repo.GetAccountVisitCounts(filters, userRole, territoryID, userID)
}

// func (s *accountService) ProcessAccount(data []string) error {
// 	existingAccount, err := s.repo.FindByAccountCode(data[8]) // AccountCode
// 	if err != nil {
// 		return err
// 	}
// 	if existingAccount != nil {
// 		return nil
// 	}

// 	// Assuming you have a cityService with a method FindCityByName
// 	city, err := s.cityRepo.FindByName(data[4]) // City
// 	if err != nil {
// 		return err
// 	}
// 	if city == nil {
// 		return fmt.Errorf("city not found")
// 	}

// 	account := models.Account{
// 		AccountName:             &data[5],
// 		AccountType:             &data[6],
// 		AccountCategory:         &data[7],
// 		AccountCode:             &data[8],
// 		City:                    &city.ID, // Set city ID
// 		ContactName:             &data[9],
// 		EmailAccount:            &data[10],
// 		WebsiteAccount:          &data[12],
// 		Potensi:                 &data[11],
// 		SystemInformasiAkademik: &data[13],
// 		Ownership:               &data[14],
// 	}
// 	return s.repo.Create(&account)
// }

func (s *accountService) ProcessAccount(data []string) error {
	if isZeroValue(data[14]) || isZeroValue(data[15]) {
		fmt.Println("data not found")
		return nil
	} else {
		fmt.Println("data found")

		existingAccount, err := s.repo.FindByAccountCode(data[0]) // AccountCode
		if err != nil {
			return err
		}

		if existingAccount == nil {
			fmt.Println("account not found")
			return nil
		}

		updateData := map[string]interface{}{
			"longitude": data[15],
			"latitude":  data[14],
		}

		return s.repo.UpdateFields(existingAccount.ID, updateData)
	}

	return nil

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
		City:                    &city.ID, // Set city ID as *uint
		ContactName:             &data[9],
		EmailAccount:            &data[10],
		WebsiteAccount:          &data[12],
		Potensi:                 &data[11],
		SystemInformasiAkademik: &data[13],
		Ownership:               &data[14],
	}
	return s.repo.Create(&account)
}

// Helper function to convert a string to a pointer
func stringPointer(s string) *string {
	return &s
}

func isZeroValue(value string) bool {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return parsed == 0
}

func (s *accountService) FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*models.Account, error) {
	account, err := s.repo.FindByAccountID(id, userRole, territoryID, userID)
	if err != nil {
		return nil, err
	}
	return account, nil
}
