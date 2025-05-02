package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

type ValidateCreateRequest struct {
	CommunicationType   string  `json:"communication_type" validate:"required"`
	Note                string  `json:"note" validate:"required"`
	AccountID           *string `json:"account_id" validate:"required"`
	ContactID           *string `json:"contact_id"`
	CheckOpportunity    *string `json:"check_opportunity" validate:"omitempty,oneof=0 1"`
	OpportunityName     *string `json:"opportunity_name"`
	StatusCommunication *string `json:"status_communication"`
}

var validationMessages = map[string]string{
	"account_id.required":         "Account ID wajib diisi",
	"communication_type.required": "Tipe komunikasi wajib diisi",
	"note.required":               "Catatan wajib diisi",
	"check_opportunity.oneof":     "Hanya bisa memilih 0 atau 1",
}

var validate = validator.New()

func ValidateCreate(req *ValidateCreateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateStatus(req *ValidateCreateRequest) map[string]string {
	errors := make(map[string]string)

	// Validate StatusCommunication
	if req.StatusCommunication == nil && req.CommunicationType == "MENAWARKAN PROGRAM" {
		errors["status_communication"] = "Status komunikasi wajib diisi"
	}

	// Return nil jika tidak ada error
	if len(errors) == 0 {
		return nil
	}

	return errors
}
