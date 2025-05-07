package repository

import "byu-crm-service/modules/eligibility/response"

type EligibilityRepository interface {
	GetEligibilities(subjectType string, categories []string, types []string, locationFilter response.LocationFilter) ([]response.Eligibility, error)
}
