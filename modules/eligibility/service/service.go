package service

type EligibilityService interface {
	CreateEligibility(subjectType string, subjectID uint, categories []string, types []string, locations map[string][]string) error
}
