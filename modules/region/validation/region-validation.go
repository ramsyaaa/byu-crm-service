package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateRegionRequest struct for creating a Area
type CreateRegionRequest struct {
	Name   string `json:"name" validate:"required"`
	AreaID int    `json:"area_id" validate:"required"`
}

type UpdateRegionRequest struct {
	Name   string `json:"name" validate:"required"`
	AreaID int    `json:"area_id" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":    "Nama Area harus diisi",
	"area_id.required": "Area harus di pilih",
}

// ValidateCreate function to validate CreateRegionRequest
func ValidateCreate(req *CreateRegionRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateRegionRequest
func ValidateUpdate(req *UpdateRegionRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
