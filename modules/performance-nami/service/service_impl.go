package service

import (
	"byu-crm-service/models"
	AccountRepository "byu-crm-service/modules/account/repository"
	cityRepository "byu-crm-service/modules/city/repository"
	"byu-crm-service/modules/performance-nami/repository"
	"fmt"
	"strings"
	"time"
)

type performanceNamiService struct {
	repo        repository.PerformanceNamiRepository
	cityRepo    cityRepository.CityRepository
	accountRepo AccountRepository.AccountRepository
}

func NewPerformanceNamiService(repo repository.PerformanceNamiRepository, cityRepo cityRepository.CityRepository, accountRepo AccountRepository.AccountRepository) PerformanceNamiService {
	return &performanceNamiService{repo: repo, cityRepo: cityRepo, accountRepo: accountRepo}
}

func (s *performanceNamiService) ProcessPerformanceNami(data []string) error {
	// Assuming you have a cityService with a method FindCityByName
	city, err := s.cityRepo.GetCityByName(data[15]) // City
	if err != nil {
		return err
	}
	if city == nil {
		return fmt.Errorf("city not found")
	}

	account, err := s.accountRepo.FindByAccountName(data[4]) // City
	if err != nil {
		return err
	}

	account_id := 0
	if account != nil {
		account_id = int(account.ID)
	}

	existingRecord, err := s.repo.FindBySerialNumberMsisdn(data[16])
	if err != nil {
		return fmt.Errorf("failed to check existing record: %w", err)
	}
	if existingRecord != nil {
		return nil
	}

	performanceNami := models.PerformanceNami{
		Periode:           &data[0],
		PeriodeDate:       parseDate(data[1]),
		EventID:           &data[2],
		PoiID:             &data[3],
		PoiName:           &data[4],
		PoiType:           &data[5],
		EventName:         &data[6],
		EventType:         &data[7],
		EventLocationType: &data[8],
		SalesType:         &data[9],
		SalesType2:        &data[10],
		// Area:                &data[11],
		// Regional:            &data[12],
		// Branch:              &data[13],
		// Cluster:             &data[14],
		CityID:             city.ID,
		SerialNumberMsisdn: &data[16],
		ScanType:           &data[17],
		ActiveMsisdn:       &data[18],
		ActiveDate:         parseDate(data[19]),
		ActiveCity:         &data[20],
		Validation:         &data[21],
		ValidKpi:           parseBool(data[22]), // Assuming parseBool is a helper to parse boolean values
		Revenue:            &data[23],
		SaDate:             parseDate(data[24]),
		SoDate:             parseDate(data[25]),
		NewImei:            dereferenceString(parseNewImei(data[26])), // Assuming parseUint8 is a helper to parse uint8 values
		SkulIDDate:         parseDate(data[27]),
		AgentID:            &data[28],
		UserID:             &data[29],
		UserName:           &data[30],
		UserType:           &data[31],
		UserSubType:        &data[32],
		ScanDate:           parseDateTime(data[33]),
		Plan:               &data[34],
		TopStatus:          parseBool(data[35]),
		AccountID:          uint(account_id),
	}
	return s.repo.Create(&performanceNami)
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

func parseNewImei(imeiStr string) *string {
	if imeiStr == "\\N" || strings.TrimSpace(imeiStr) == "" {
		return nil
	}

	if imeiStr == "-" {
		str := "0"
		return &str
	}

	if imeiStr == "N" {
		str := "1"
		return &str
	}

	if imeiStr == "Y" {
		str := "2"
		return &str
	}

	return nil
}

func dereferenceString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func parseDateTime(dateTimeStr string) *time.Time {
	if dateTimeStr == "\\N" || strings.TrimSpace(dateTimeStr) == "" {
		return nil
	}

	parsedDateTime, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		fmt.Printf("Error parsing datetime: %s\n", err)
		return nil
	}
	return &parsedDateTime
}

func parseBool(value string) bool {
	value = strings.TrimSpace(value)
	return value == "Y"
}
