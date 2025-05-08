package service

import (
	"byu-crm-service/modules/eligibility/repository"
)

type eligibilityService struct {
	repo repository.EligibilityRepository
}

func NewEligibilityService(repo repository.EligibilityRepository) EligibilityService {
	return &eligibilityService{repo: repo}
}

func (s *eligibilityService) CreateEligibility(subjectType string, subjectID uint, categories []string, types []string, locations map[string][]string) error {

	err := s.repo.InsertEligibility(
		subjectType,
		subjectID,
		categories,
		types,
		locations,
	)

	return err
}
