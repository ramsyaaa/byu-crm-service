package service

type PerformanceIndianaService interface {
	ProcessPerformanceIndiana(data []string, month uint, year uint) error
}
