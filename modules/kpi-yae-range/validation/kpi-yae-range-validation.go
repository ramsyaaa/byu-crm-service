package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateKpiYaeRangeRequest struct for creating a KPI YAE range
type CreateKpiYaeRangeRequest struct {
	StartDate string   `json:"start_date" validate:"required"`
	EndDate   string   `json:"end_date" validate:"required"`
	Name      []string `json:"name" validate:"required,min=1,dive,required"`
	Target    []string `json:"target" validate:"required,dive"`
}

var validationMessages = map[string]string{
	"StartDate.required": "Tanggal mulai harus diisi",
	"EndDate.required":   "Tanggal selesai harus diisi",
	"Name.required":      "Nama KPI harus diisi",
	"Name.min":           "Minimal satu nama KPI harus diisi",
	"Name[].required":    "Semua nama KPI tidak boleh kosong",
	"Target.required":    "Target KPI harus diisi",
}

// ValidateCreate function to validate CreateKpiYaeRangeRequest
func ValidateCreate(req *CreateKpiYaeRangeRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
