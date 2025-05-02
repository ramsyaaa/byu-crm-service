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
	"opportunity_name.required": "Nama opportunity wajib diisi",
	"description.required":      "Deskripsi wajib diisi",
}

// ValidateData function to validate ValidateRequest
func ValidateData(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
