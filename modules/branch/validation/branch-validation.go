package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateBranchRequest struct for creating a Area
type CreateBranchRequest struct {
	Name     string `json:"name" validate:"required"`
	RegionID int    `json:"region_id" validate:"required"`
}

type UpdateBranchRequest struct {
	Name     string `json:"name" validate:"required"`
	RegionID int    `json:"region_id" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":      "Nama Area harus diisi",
	"region_id.required": "Regional harus di pilih",
}

// ValidateCreate function to validate CreateBranchRequest
func ValidateCreate(req *CreateBranchRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateBranchRequest
func ValidateUpdate(req *UpdateBranchRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
