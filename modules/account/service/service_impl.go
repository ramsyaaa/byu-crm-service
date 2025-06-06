package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/account/response"
	cityRepository "byu-crm-service/modules/city/repository"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type accountService struct {
	repo     repository.AccountRepository
	cityRepo cityRepository.CityRepository
}

func NewAccountService(repo repository.AccountRepository, cityRepo cityRepository.CityRepository) AccountService {
	return &accountService{repo: repo, cityRepo: cityRepo}
}

func (s *accountService) GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool, excludeVisited bool) ([]response.AccountResponse, int64, error) {
	return s.repo.GetAllAccounts(limit, paginate, page, filters, userRole, territoryID, userID, onlyUserPic, excludeVisited)
}

func (s *accountService) CreateAccount(requestBody map[string]interface{}, userID int) ([]models.Account, error) {
	// Use getStringValue to safely handle nil values and type conversions
	accountData := map[string]string{
		"account_name":              getStringValue(requestBody["account_name"]),
		"account_type":              getStringValue(requestBody["account_type"]),
		"account_category":          getStringValue(requestBody["account_category"]),
		"account_code":              getStringValue(requestBody["account_code"]),
		"city":                      getStringValue(requestBody["city"]),
		"contact_name":              getStringValue(requestBody["contact_name"]),
		"email_account":             getStringValue(requestBody["email_account"]),
		"website_account":           getStringValue(requestBody["website_account"]),
		"system_informasi_akademik": getStringValue(requestBody["system_informasi_akademik"]),
		"ownership":                 getStringValue(requestBody["ownership"]),
		"pic":                       getStringValue(requestBody["pic"]),
		"pic_internal":              getStringValue(requestBody["pic_internal"]),
		"latitude":                  getStringValue(requestBody["latitude"]),
		"longitude":                 getStringValue(requestBody["longitude"]),
	}

	accounts, err := s.repo.CreateAccount(accountData, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountService) UpdateAccount(requestBody map[string]interface{}, accountID int, userRole string, territoryID int, userID int) ([]models.Account, error) {
	accountData := map[string]string{
		"account_name":              getStringValue(requestBody["account_name"]),
		"account_type":              getStringValue(requestBody["account_type"]),
		"account_category":          getStringValue(requestBody["account_category"]),
		"account_code":              getStringValue(requestBody["account_code"]),
		"city":                      getStringValue(requestBody["city"]),
		"contact_name":              getStringValue(requestBody["contact_name"]),
		"email_account":             getStringValue(requestBody["email_account"]),
		"website_account":           getStringValue(requestBody["website_account"]),
		"system_informasi_akademik": getStringValue(requestBody["system_informasi_akademik"]),
		"ownership":                 getStringValue(requestBody["ownership"]),
		"pic":                       getStringValue(requestBody["pic"]),
		"pic_internal":              getStringValue(requestBody["pic_internal"]),
		"latitude":                  getStringValue(requestBody["latitude"]),
		"longitude":                 getStringValue(requestBody["longitude"]),
	}

	accounts, err := s.repo.UpdateAccount(accountData, accountID, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountService) UpdateFields(id uint, fields map[string]interface{}) error {
	return s.repo.UpdateFields(id, fields)
}

func (s *accountService) UpdatePic(accountID int, userRole string, territoryID int, userID int) (*response.AccountResponse, error) {
	existingAccount, err := s.repo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))

	if err != nil {
		return nil, err
	}

	if existingAccount.Pic != nil && *existingAccount.Pic != "" {
		return nil, errors.New("PIC sudah diatur dan tidak dapat diubah")
	}

	err = s.repo.UpdateFields(uint(accountID), map[string]interface{}{"pic": userID})
	if err != nil {
		return nil, err
	}

	existingAccount, err = s.repo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))

	if err != nil {
		return nil, err
	}

	return existingAccount, nil
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
	city, err := s.cityRepo.GetCityByName(data[4]) // City
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

func (s *accountService) CheckAlreadyUpdateData(accountID int, clockIn time.Time, userID int) (bool, error) {
	return s.repo.CheckAlreadyUpdateData(accountID, userID, clockIn)
}

func (s *accountService) CreateHistoryActivityAccount(userID, accountID uint, updateType string, subjectType *string, subjectID *uint) error {
	return s.repo.CreateHistoryActivityAccount(userID, accountID, updateType, subjectType, subjectID)
}

func isZeroValue(value string) bool {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return parsed == 0
}

func (s *accountService) FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.SingleAccountResponse, error) {
	account, err := s.repo.FindByAccountID(id, userRole, territoryID, userID)
	if err != nil {
		return nil, err
	}

	var accountResponse response.SingleAccountResponse

	accountResponse.ID = account.ID
	accountResponse.AccountName = account.AccountName
	accountResponse.AccountType = account.AccountType
	accountResponse.AccountCategory = account.AccountCategory
	accountResponse.AccountCode = account.AccountCode
	accountResponse.City = account.City
	accountResponse.CityName = account.CityName
	accountResponse.ClusterName = account.ClusterName
	accountResponse.BranchName = account.BranchName
	accountResponse.RegionName = account.RegionName
	accountResponse.AreaName = account.AreaName
	accountResponse.ContactName = account.ContactName
	accountResponse.EmailAccount = account.EmailAccount
	accountResponse.WebsiteAccount = account.WebsiteAccount
	accountResponse.Potensi = account.Potensi
	accountResponse.SystemInformasiAkademik = account.SystemInformasiAkademik
	accountResponse.CustomerSegmentationId = account.CustomerSegmentationId
	accountResponse.Latitude = account.Latitude
	accountResponse.Longitude = account.Longitude
	accountResponse.Ownership = account.Ownership
	accountResponse.Pic = account.Pic
	accountResponse.PicInternal = account.PicInternal
	accountResponse.SocialMedias = account.SocialMedias
	accountResponse.Contacts = account.Contacts
	accountResponse.Products = account.Products
	if account.PicDetail != nil {
		accountResponse.PicDetail = &models.UserResponse{
			ID:    account.PicDetail.ID,
			Name:  account.PicDetail.Name,
			Email: account.PicDetail.Email,
		}
	} else {
		accountResponse.PicDetail = nil
	}
	if account.PicInternalDetail != nil {
		accountResponse.PicInternalDetail = &models.UserResponse{
			ID:    account.PicInternalDetail.ID,
			Name:  account.PicInternalDetail.Name,
			Email: account.PicInternalDetail.Email,
		}
	} else {
		accountResponse.PicInternalDetail = nil
	}
	accountResponse.IsSkulid = account.IsSkulid
	accountResponse.CreatedAt = account.CreatedAt
	accountResponse.UpdatedAt = account.UpdatedAt

	accountResponse.Category = []string{}
	accountResponse.Url = []string{}
	accountResponse.ProductAccount = []string{}

	for _, sm := range accountResponse.SocialMedias {
		accountResponse.Category = append(accountResponse.Category, *sm.Category)
		accountResponse.Url = append(accountResponse.Url, *sm.Url)
	}

	for _, sm := range accountResponse.Products {
		accountResponse.ProductAccount = append(accountResponse.ProductAccount, fmt.Sprintf("%d", sm.ID))
	}

	if len(account.Contacts) > 0 {
		var contactIDs []string
		for _, contact := range account.Contacts {
			if contact.ID != 0 {
				contactIDs = append(contactIDs, fmt.Sprintf("%d", contact.ID))
			}
		}
		accountResponse.ContactID = contactIDs
	}

	if *account.AccountCategory == "SEKOLAH" {
		if account.AccountTypeSchoolDetail != nil {
			accountResponse.DiesNatalis = &account.AccountTypeSchoolDetail.DiesNatalis
			accountResponse.Extracurricular = account.AccountTypeSchoolDetail.Extracurricular
			accountResponse.FootballFieldBranding = account.AccountTypeSchoolDetail.FootballFieldBrannnding
			accountResponse.BasketballFieldBranding = account.AccountTypeSchoolDetail.BasketballFieldBranding
			accountResponse.WallPaintingBranding = account.AccountTypeSchoolDetail.WallPaintingBranding
			accountResponse.WallMagazineBranding = account.AccountTypeSchoolDetail.WallMagazineBranding
		}

	} else if *account.AccountCategory == "KAMPUS" {
		if account.AccountTypeCampusDetail != nil {
			if account.AccountTypeCampusDetail.RangeAge != nil {
				rangeAge, err := SplitFields(*account.AccountTypeCampusDetail.RangeAge, []string{"age", "percentage_age"})
				if err != nil {
					accountResponse.Age = []string{}
					accountResponse.PercentageAge = []string{}
				} else {
					accountResponse.Age = rangeAge["age"]
					accountResponse.PercentageAge = rangeAge["percentage_age"]
				}
			}

			if account.AccountTypeCampusDetail.Origin != nil {
				origin, err := SplitFields(*account.AccountTypeCampusDetail.Origin, []string{"origin", "percentage_origin"})
				if err != nil {
					accountResponse.Origin = []string{}
					accountResponse.PercentageOrigin = []string{}
				} else {
					accountResponse.Origin = origin["origin"]
					accountResponse.PercentageOrigin = origin["percentage_origin"]
				}
			}

			if account.AccountTypeCampusDetail.OrganizationName != nil {
				organization_name, _ := SplitFields(*account.AccountTypeCampusDetail.OrganizationName, []string{"organization_name"})
				if err != nil {
					accountResponse.OrganizationName = []string{}
				} else {
					accountResponse.OrganizationName = organization_name["organization_name"]
				}
			}

			if account.AccountTypeCampusDetail.PreferenceTechnologies != nil {
				preferenceTechnologies, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.PreferenceTechnologies)
				if err != nil {
					accountResponse.PreferenceTechnologies = []string{}
				} else {
					accountResponse.PreferenceTechnologies = preferenceTechnologies
				}
			}

			if account.AccountTypeCampusDetail.MemberNeeds != nil {
				memberNeeds, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.MemberNeeds)
				if err != nil {
					accountResponse.MemberNeeds = []string{}
				} else {
					accountResponse.MemberNeeds = memberNeeds
				}
			}

			if account.AccountTypeCampusDetail.ItInfrastructures != nil {
				itInfrastructures, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.ItInfrastructures)
				if err != nil {
					accountResponse.ItInfrastructures = []string{}
				} else {
					accountResponse.ItInfrastructures = itInfrastructures
				}
			}

			if account.AccountTypeCampusDetail.DigitalCollaborations != nil {
				digitalCollaborations, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.DigitalCollaborations)
				if err != nil {
					accountResponse.DigitalCollaborations = []string{}
				} else {
					accountResponse.DigitalCollaborations = digitalCollaborations
				}
			}
			accountResponse.AccessTechnology = account.AccountTypeCampusDetail.AccessTechnology
			accountResponse.CampusAdministrationApp = account.AccountTypeCampusDetail.CampusAdministrationApp

			accountResponse.PotentionalCollaboration = account.AccountTypeCampusDetail.PotentionalCollaboration

			if account.AccountTypeCampusDetail.UniversityRank != nil {
				rank, err := SplitFields(*account.AccountTypeCampusDetail.UniversityRank, []string{"rank", "year_rank"})
				if err != nil {
					accountResponse.Rank = []string{}
					accountResponse.YearRank = []string{}
				} else {
					accountResponse.Rank = rank["rank"]
					accountResponse.YearRank = rank["year_rank"]
				}
			}

			if account.AccountTypeCampusDetail.FocusProgramStudy != nil {
				programStudy, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.FocusProgramStudy)
				if err != nil {
					accountResponse.ProgramStudy = []string{}
				} else {
					accountResponse.ProgramStudy = programStudy
				}
			}

			if account.AccountTypeCampusDetail.ProgramIdentification != nil {
				programIdentification, err := parseJSONStringToArray(*account.AccountTypeCampusDetail.ProgramIdentification)
				if err != nil {
					accountResponse.ProgramIdentification = []string{}
				} else {
					accountResponse.ProgramIdentification = programIdentification
				}
			}

			if account.AccountTypeCampusDetail.Byod != nil {
				byodStr := fmt.Sprintf("%d", *account.AccountTypeCampusDetail.Byod)
				accountResponse.Byod = &byodStr
			} else {
				accountResponse.Byod = nil
			}
		}

		var years []string
		var amounts []string

		if len(account.AccountMembers) > 0 {
			for _, member := range account.AccountMembers {
				if member.Year != nil {
					years = append(years, *member.Year)
				}
				if member.Amount != nil {
					amounts = append(amounts, *member.Amount)
				}
			}
		}

		accountResponse.Year = years
		accountResponse.Amount = amounts
		years = []string{}
		amounts = []string{}

		if len(account.AccountLectures) > 0 {
			for _, member := range account.AccountLectures {
				if member.Year != nil {
					years = append(years, *member.Year)
				}
				if member.Amount != nil {
					amounts = append(amounts, *member.Amount)
				}
			}
		}

		accountResponse.YearLecture = years
		accountResponse.AmountLecture = amounts

		faculties := []string{}
		for _, faculty := range account.AccountFaculties {
			if faculty.FacultyID != nil {
				faculties = append(faculties, fmt.Sprintf("%d", *faculty.FacultyID))
			}
		}

		accountResponse.Faculties = faculties
	} else if *account.AccountCategory == "KOMUNITAS" {

		if account.AccountTypeCommunityDetail != nil {
			if account.AccountTypeCommunityDetail.PotentionalCollaboration != nil {
				accountResponse.PotentionalCollaboration = account.AccountTypeCommunityDetail.PotentionalCollaboration
			} else {
				accountResponse.PotentionalCollaboration = nil
			}

			if account.AccountTypeCommunityDetail.AccountSubtype != nil {
				accountResponse.AccountSubtype = account.AccountTypeCommunityDetail.AccountSubtype
			} else {
				accountResponse.AccountSubtype = nil
			}

			if account.AccountTypeCommunityDetail.Group != nil {
				accountResponse.Group = account.AccountTypeCommunityDetail.Group
			} else {
				accountResponse.Group = nil
			}

			if account.AccountTypeCommunityDetail.GroupName != nil {
				accountResponse.GroupName = account.AccountTypeCommunityDetail.GroupName
			} else {
				accountResponse.GroupName = nil
			}

			if account.AccountTypeCommunityDetail.ProductService != nil {
				accountResponse.ProductService = account.AccountTypeCommunityDetail.ProductService
			} else {
				accountResponse.ProductService = nil
			}

			if account.AccountTypeCommunityDetail.PotentialCollaborationItems != nil {
				accountResponse.PotentialCollaborationItems = account.AccountTypeCommunityDetail.PotentialCollaborationItems
			} else {
				accountResponse.PotentialCollaborationItems = nil
			}
		} else {
			accountResponse.PotentionalCollaboration = nil
			accountResponse.AccountSubtype = nil
			accountResponse.Group = nil
			accountResponse.GroupName = nil
			accountResponse.ProductService = nil
			accountResponse.PotentialCollaborationItems = nil
		}

		var years []string
		var amounts []string

		if len(account.AccountMembers) > 0 {
			for _, member := range account.AccountMembers {
				if member.Year != nil {
					years = append(years, *member.Year)
				}
				if member.Amount != nil {
					amounts = append(amounts, *member.Amount)
				}
			}
		}

		accountResponse.Year = years
		accountResponse.Amount = amounts

		// Cek data AccountTypeCampusDetail
		if account.AccountTypeCampusDetail != nil {
			if account.AccountTypeCampusDetail.RangeAge != nil {
				rangeAge, err := SplitFields(*account.AccountTypeCampusDetail.RangeAge, []string{"age", "percentage_age"})
				if err != nil {
					accountResponse.Age = []string{}
					accountResponse.PercentageAge = []string{}
				} else {
					accountResponse.Age = rangeAge["age"]
					accountResponse.PercentageAge = rangeAge["percentage_age"]
				}
			}
		}

		// Cek data AccountTypeCommunityDetail
		if account.AccountTypeCommunityDetail != nil {
			if account.AccountTypeCommunityDetail.Gender != nil {
				gender, err := SplitFields(*account.AccountTypeCommunityDetail.Gender, []string{"gender", "percentage_gender"})
				if err != nil {
					accountResponse.Gender = []string{}
					accountResponse.PercentageGender = []string{}
				} else {
					accountResponse.Gender = gender["gender"]
					accountResponse.PercentageGender = gender["percentage_gender"]
				}
			}

			if account.AccountTypeCommunityDetail.EducationalBackground != nil {
				educationalBackground, err := SplitFields(*account.AccountTypeCommunityDetail.EducationalBackground, []string{"educational_background", "percentage_educational_background"})
				if err != nil {
					accountResponse.EducationalBackground = []string{}
					accountResponse.PercentageEducationalBackground = []string{}
				} else {
					accountResponse.EducationalBackground = educationalBackground["educational_background"]
					accountResponse.PercentageEducationalBackground = educationalBackground["percentage_educational_background"]
				}
			}

			if account.AccountTypeCommunityDetail.Profession != nil {
				profession, err := SplitFields(*account.AccountTypeCommunityDetail.Profession, []string{"profession", "percentage_profession"})
				if err != nil {
					accountResponse.Profession = []string{}
					accountResponse.PercentageProfession = []string{}
				} else {
					accountResponse.Profession = profession["profession"]
					accountResponse.PercentageProfession = profession["percentage_profession"]
				}
			}

			if account.AccountTypeCommunityDetail.Income != nil {
				income, err := SplitFields(*account.AccountTypeCommunityDetail.Income, []string{"income", "percentage_income"})
				if err != nil {
					accountResponse.Income = []string{}
					accountResponse.PercentageIncome = []string{}
				} else {
					accountResponse.Income = income["income"]
					accountResponse.PercentageIncome = income["percentage_income"]
				}
			}
		}

	}

	return &accountResponse, nil
}

func (s *accountService) CountAccount(userRole string, territoryID int, withGeoJson bool) (int64, map[string]int64, []map[string]interface{}, response.TerritoryInfo, error) {
	count, categories, territories, territory_info, err := s.repo.CountAccount(userRole, territoryID, withGeoJson)
	if err != nil {
		return 0, nil, nil, response.TerritoryInfo{}, err
	}
	return count, categories, territories, territory_info, nil
}

// FindAccountsWithDifferentPic retrieves accounts that have a different PIC than the specified user ID.
func (s *accountService) FindAccountsWithDifferentPic(accountIDs []int, userID int) ([]models.Account, error) {
	// Validasi jika accountIDs kosong
	if len(accountIDs) == 0 {
		return nil, errors.New("no account IDs provided")
	}
	accounts, err := s.repo.FindAccountsWithDifferentPic(accountIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding accounts with different PIC: %v", err)
	}
	return accounts, nil
}

func (s *accountService) UpdatePicMultipleAccounts(accountIDs []int, userID int) error {
	// Validasi jika accountIDs kosong
	if len(accountIDs) == 0 {
		return errors.New("no account IDs provided")
	}

	err := s.repo.UpdatePicMultipleAccounts(accountIDs, userID)
	if err != nil {
		return fmt.Errorf("error updating PIC for account ID %d: %v", accountIDs, err)
	}

	return nil
}

func SplitFields(jsonString string, keys []string) (map[string][]string, error) {
	// Cek jika JSON string kosong atau null
	if jsonString == "" || jsonString == "null" {
		// Jika kosong atau null, kembalikan array kosong untuk setiap key
		result := make(map[string][]string)
		for _, key := range keys {
			result[key] = []string{}
		}
		return result, nil
	}

	// Slice untuk menampung hasil decoding JSON
	var rawData []map[string]string

	// Decode JSON menjadi slice of map[string]string
	err := json.Unmarshal([]byte(jsonString), &rawData)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// Map untuk menyimpan hasil pemisahan berdasarkan key
	result := make(map[string][]string)

	// Iterasi untuk setiap record dalam rawData
	for _, item := range rawData {
		// Iterasi setiap key yang ingin dipisahkan
		for _, key := range keys {
			// Cek apakah key ada dalam item, jika ada, tambahkan ke hasil
			if value, exists := item[key]; exists {
				result[key] = append(result[key], value)
			} else {
				// Jika key tidak ada, tambahkan array kosong
				result[key] = append(result[key], "")
			}
		}
	}

	return result, nil
}

func parseJSONStringToArray(jsonString string) ([]string, error) {
	// Jika jsonString kosong atau null, kembalikan array kosong
	if jsonString == "" || jsonString == "null" {
		return []string{}, nil
	}

	// Slice untuk menyimpan hasil decoding JSON
	var result []string

	// Decode JSON menjadi slice of string
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return result, nil
}

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
