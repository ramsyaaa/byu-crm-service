package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/performance-indiana/repository"
	userRepository "byu-crm-service/modules/user/repository"
	"fmt"
	"strconv"
)

type performanceIndianaService struct {
	repo     repository.PerformanceIndianaRepository
	userRepo userRepository.UserRepository
}

func NewPerformanceIndianaService(repo repository.PerformanceIndianaRepository, userRepo userRepository.UserRepository) PerformanceIndianaService {
	return &performanceIndianaService{repo: repo, userRepo: userRepo}
}

func (s *performanceIndianaService) ProcessPerformanceIndiana(
	data []string,
	month uint,
	year uint,
) error {
	user, err := s.userRepo.FindByYaeCode(data[3])
	if err != nil {
		return err
	}

	userID := int(user.ID)

	dataIn, err := ParseStringToInt(data[5])
	if err != nil {
		return err
	}
	notProcess, err := ParseStringToInt(data[6])
	if err != nil {
		return err
	}
	rejected, err := ParseStringToInt(data[7])
	if err != nil {
		return err
	}
	pending, err := ParseStringToInt(data[8])
	if err != nil {
		return err
	}
	approve, err := ParseStringToInt(data[9])
	if err != nil {
		return err
	}
	inProgress, err := ParseStringToInt(data[10])
	if err != nil {
		return err
	}
	failed, err := ParseStringToInt(data[11])
	if err != nil {
		return err
	}
	active, err := ParseStringToInt(data[12])
	if err != nil {
		return err
	}
	activeBySystem, err := ParseStringToInt(data[13])
	if err != nil {
		return err
	}

	// 🔑 month format: YYYY-MM (konsisten dgn input <input type="month">)
	monthString := fmt.Sprintf("%04d-%02d", year, month)

	existing, err := s.repo.FindByUserAndMonth(userID, monthString)
	if err != nil {
		return err
	}

	performance := models.PerformanceIndiana{
		UserID:         &userID,
		DataIn:         &dataIn,
		NotProcess:     &notProcess,
		Rejected:       &rejected,
		Pending:        &pending,
		Approve:        &approve,
		InProgress:     &inProgress,
		Failed:         &failed,
		Active:         &active,
		ActiveBySystem: &activeBySystem,
		Month:          &monthString,
	}

	// 🔁 UPDATE jika ada
	if existing != nil {
		performance.ID = existing.ID
		return s.repo.Update(&performance)
	}

	// ➕ CREATE jika belum ada
	return s.repo.Create(&performance)
}

func (s *performanceIndianaService) GetDataInByUserAndMonth(
	userID int,
	month uint,
	year uint,
) (int, error) {
	return s.repo.GetDataInByUserAndMonth(userID, month, year)
}

func ParseStringToInt(stringData string) (int, error) {
	dataIn, err := strconv.ParseInt(stringData, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(dataIn), nil
}
