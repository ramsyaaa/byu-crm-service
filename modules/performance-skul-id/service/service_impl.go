package service

import (
	"byu-crm-service/models"
	accountRepo "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/performance-skul-id/repository"
	subdistrictRepo "byu-crm-service/modules/subdistrict/repository"
	"fmt"
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

	performanceSkulId := models.PerformanceSkulId{
		UserId:         account.Pic,
		UserType:       &data[1],
		RegisteredDate: parseDate(data[2]),
		Msisdn:         &data[3],
		Provider:       &data[4],
		AccountId:      &account.ID,
		UserName:       &data[6],
		FlagNewSales:   &data[18],
		FlagImei:       &data[19],
		RevMtd:         &data[20],
		RevMtdM1:       &data[21],
		RevDigital:     &data[22],
		ActivityMtd:    &data[23],
		FlagActiveMtd:  &data[24],
		SubdistrictId:  &subdistrict.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.repo.Create(&performanceSkulId)
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
