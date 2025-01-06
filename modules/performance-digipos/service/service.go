package service

type PerformanceDigiposService interface {
	ProcessPerformanceDigipos(data []string) error
}
