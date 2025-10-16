package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateCalculatorKpiRequest struct for creating a KPI YAE range
type CreateCalculatorKpiRequest struct {
	StartDate string   `json:"start_date" validate:"required"`
	EndDate   string   `json:"end_date" validate:"required"`
	Incentive80      string `json:"incentive_80" validate:"required"`
	Incentive90      string `json:"incentive_90" validate:"required"`
	Incentive100      string `json:"incentive_100" validate:"required"`
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

// ValidateCreate function to validate CreateCalculatorKpiRequest
func ValidateCreate(req *CreateCalculatorKpiRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
