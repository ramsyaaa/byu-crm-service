package service

type PerformanceIndianaService interface {
	ProcessPerformanceIndiana(data []string, month uint, year uint) error
	GetDataInByUserAndMonth(userID int, month uint, year uint) (int, error)
}
