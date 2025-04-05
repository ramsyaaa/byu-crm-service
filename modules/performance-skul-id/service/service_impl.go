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

	subdistrict, err := s.subdistrictRepo.FindByName(data[13])
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

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Printf("Error parsing date: %s\n", err)
		return nil
	}
	return &parsedDate
}
