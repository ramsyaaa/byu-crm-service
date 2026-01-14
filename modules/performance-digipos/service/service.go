package service

type PerformanceDigiposService interface {
	ProcessPerformanceDigipos(data []string) error
	CountPerformanceByUserYaeCode(user_id int, month uint, year uint) (int, error)
}
