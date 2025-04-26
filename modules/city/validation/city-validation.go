package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateCityRequest struct for creating a Area
type CreateCityRequest struct {
	Name      string `json:"name" validate:"required"`
	ClusterID int    `json:"cluster_id" validate:"required"`
}

type UpdateCityRequest struct {
	Name      string `json:"name" validate:"required"`
	ClusterID int    `json:"cluster_id" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":       "Nama Area harus diisi",
	"cluster_id.required": "Cluster harus di pilih",
}

// ValidateCreate function to validate CreateCityRequest
func ValidateCreate(req *CreateCityRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateCityRequest
func ValidateUpdate(req *UpdateCityRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
