package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// ValidateRequest struct for creating a Area
type ValidateRequest struct {
	OpportunityName string  `json:"opportunity_name" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	OpenDate        *string `json:"open_date"`
	CloseDate       *string `json:"close_date"`
	AccountID       *string `json:"account_id"`
	ContactID       *string `json:"contact_id"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"OpportunityName.required": "Opportunity name is required",
	"Description.required":     "Description is required",
}

// ValidateCreate function to validate ValidateRequest
func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
