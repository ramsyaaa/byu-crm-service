package repository

import "byu-crm-service/modules/eligibility/response"

type EligibilityRepository interface {
	GetEligibilities(subjectType string, categories []string, types []string, locationFilter response.LocationFilter) ([]response.Eligibility, error)
	InsertEligibility(subjectType string, subjectID uint, categories []string, types []string, locations map[string][]string) error
}
