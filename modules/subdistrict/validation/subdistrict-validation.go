package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateSubdistrictRequest struct for creating a Area
type CreateSubdistrictRequest struct {
	Name   string `json:"name" validate:"required"`
	CityID int    `json:"city_id" validate:"required"`
}

type UpdateSubdistrictRequest struct {
	Name   string `json:"name" validate:"required"`
	CityID int    `json:"city_id" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":    "Nama Area harus diisi",
	"city_id.required": "City harus di pilih",
}

// ValidateCreate function to validate CreateSubdistrictRequest
func ValidateCreate(req *CreateSubdistrictRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateSubdistrictRequest
func ValidateUpdate(req *UpdateSubdistrictRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
