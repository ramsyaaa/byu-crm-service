package service

import (
	"byu-crm-service/models"
	accountRepo "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/performance-skul-id/repository"
	subdistrictRepo "byu-crm-service/modules/subdistrict/repository"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type performanceSkulIdService struct {
	repo            repository.PerformanceSkulIdRepository
	accountRepo     accountRepo.AccountRepository
	subdistrictRepo subdistrictRepo.SubdistrictRepository
}

func NewPerformanceSkulIdService(repo repository.PerformanceSkulIdRepository, accountRepo accountRepo.AccountRepository, subdistrictRepo subdistrictRepo.SubdistrictRepository) PerformanceSkulIdService {
	return &performanceSkulIdService{repo: repo, accountRepo: accountRepo, subdistrictRepo: subdistrictRepo}
}

func (s *performanceSkulIdService) ProcessPerformanceSkulId(data []string) error {
	account, err := s.accountRepo.FindByAccountCode(data[5])
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("account not found")
	}

	subdistrict, err := s.subdistrictRepo.GetSubdistrictByName(data[13])
	if err != nil {
		return err
	}
	if subdistrict == nil {
		return fmt.Errorf("subdistrict not found")
	}

	idSkulId := &data[0]
	if idSkulId == nil || *idSkulId == "" {
		return fmt.Errorf("id_skulid tidak boleh kosong")
	}

	performanceSkulId := models.PerformanceSkulId{
		UserId:         stringToUint(account.Pic),
		IdSkulid:       idSkulId,
		UserType:       &data[1],
		RegisteredDate: parseDate(data[2]),
		Msisdn:         &data[3],
		Provider:       &data[4],
		AccountId:      &account.ID,
		UserName:       &data[6],
		FlagNewSales:   boolToInt(data[18]),
		FlagImei:       boolToInt(data[19]),
		RevMtd:         &data[20],
		RevMtdM1:       &data[21],
		RevDigital:     &data[22],
		ActivityMtd:    &data[23],
		FlagActiveMtd:  &data[24],
		SubdistrictId:  &subdistrict.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	existingPerformance, err := s.repo.FindByIdSkulId(*idSkulId)
	if err != nil {
		return err
	}

	if existingPerformance != nil {
		// Update jika id_import sudah ada
		performanceSkulId.ID = existingPerformance.ID // Gunakan ID yang sudah ada
		return s.repo.Update(&performanceSkulId)
	}
	return s.repo.Create(&performanceSkulId)
}

func (s *performanceSkulIdService) ProcessPerformanceSkulIdByAccount(data []string, userID, accountID int, userRole string, territoryID int, userType string) error {
	account, err := s.accountRepo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("account not found")
	}

	idSkulId := &data[0]
	if idSkulId == nil || *idSkulId == "" {
		return fmt.Errorf("id_skulid tidak boleh kosong")
	}

	formattedUserType := capitalizeFirst(userType)

	performanceSkulId := models.PerformanceSkulId{
		UserId:         func(u int) *uint { v := uint(u); return &v }(userID),
		IdSkulid:       idSkulId,
		UserType:       &formattedUserType,
		RegisteredDate: parseDate(data[4]),
		Msisdn:         &data[2],
		Provider:       &data[3],
		AccountId:      &account.ID,
		UserName:       &data[1],
		Batch:          &data[5],
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	existingPerformance, err := s.repo.FindByIdSkulId(*idSkulId)
	if err != nil {
		return err
	}

	if existingPerformance != nil {
		updates := map[string]interface{}{}

		updates["id_skulid"] = *idSkulId
		if data[1] != "" {
			updates["user_name"] = data[1]
		}
		if data[2] != "" {
			updates["msisdn"] = data[2]
		}
		if data[3] != "" {
			updates["provider"] = data[3]
		}
		if date := parseDate(data[4]); !date.IsZero() {
			updates["registered_date"] = date
		}
		if data[5] != "" {
			updates["batch"] = data[5]
		}

		updates["user_type"] = formattedUserType
		updates["user_id"] = userID
		updates["account_id"] = account.ID
		performanceSkulId.ID = existingPerformance.ID // Gunakan ID yang sudah ada
		fmt.Printf("Updating existing performance skul id for %s\n", *idSkulId)

		return s.repo.UpdateByFields(performanceSkulId.ID, updates)
	}
	fmt.Printf("Creating new performance skul id for %s\n", *idSkulId)
	return s.repo.Create(&performanceSkulId)
}

func (s *performanceSkulIdService) FindAll(limit, offset int, filters map[string]string, accountID int, page int, paginate bool) ([]models.PerformanceSkulId, int64, error) {
	return s.repo.FindAll(limit, offset, filters, accountID, page, paginate)
}

func capitalizeFirst(s string) string {
	s = strings.ToLower(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func boolToInt(value string) *int {
	var result int
	if value == "Y" {
		result = 1
	} else {
		result = 0
	}
	return &result // Mengembalikan pointer ke int
}

func stringToUint(value *string) *uint {
	if value == nil {
		return nil
	}
	parsedValue, err := strconv.ParseUint(*value, 10, 32)
	if err != nil {
		fmt.Printf("Error converting string to uint: %s\n", err)
		return nil
	}
	result := uint(parsedValue)
	return &result
}

func parseDate(dateStr string) *time.Time {
	if dateStr == "\\N" || strings.TrimSpace(dateStr) == "" {
		return nil
	}

	formats := []string{
		"2006-01-02",      // 2025-06-08
		"02/01/2006",      // 08/06/2025
		"02-01-2006",      // 08-06-2025
		"02 Jan 2006",     // 08 Jun 2025
		"January 2, 2006", // June 8, 2025
		"2 Jan 2006",      // 8 Jun 2025
		"2006/01/02",      // 2025/06/08
		"2006.01.02",      // 2025.06.08
		"02.01.2006",      // 08.06.2025
		"02-Jan-2006",     // 08-Jun-2025
		"02-01-06",        // 09-09-20 → 2020-09-09 (dd-MM-yy)
	}

	for _, format := range formats {
		if parsedDate, err := time.Parse(format, dateStr); err == nil {
			return &parsedDate
		}
	}

	fmt.Printf("Error: unsupported date format: %s\n", dateStr)
	return nil
}

func (s *performanceSkulIdService) FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error) {
	return s.repo.FindBySerialNumberMsisdn(serial)
}

func (s *performanceSkulIdService) FindByIdSkulId(idSkulId string) (*models.PerformanceSkulId, error) {
	return s.repo.FindByIdSkulId(idSkulId)
}

func (s *performanceSkulIdService) CreatePerformanceSkulID(acccount_id int, userName, idSkulId, msisdn string, registeredDate *time.Time, provider *string, batch *string, user_type *string) (*models.PerformanceSkulId, error) {
	uAccountId := uint(acccount_id)
	formattedUserType := capitalizeFirst(*user_type)
	performance := &models.PerformanceSkulId{
		UserName:       &userName,
		UserType:       &formattedUserType,
		IdSkulid:       &idSkulId,
		Msisdn:         &msisdn,
		RegisteredDate: registeredDate,
		Provider:       provider,
		Batch:          batch,
		AccountId:      &uAccountId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := s.repo.Create(performance)
	if err != nil {
		return nil, fmt.Errorf("failed to create performance skul id: %w", err)
	}

	return performance, nil
}
