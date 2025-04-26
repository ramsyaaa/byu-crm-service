package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateClusterRequest struct for creating a Area
type CreateClusterRequest struct {
	Name     string `json:"name" validate:"required"`
	BranchID int    `json:"branch_id" validate:"required"`
}

type UpdateClusterRequest struct {
	Name     string `json:"name" validate:"required"`
	BranchID int    `json:"branch_id" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":      "Nama Area harus diisi",
	"branch_id.required": "Branch harus di pilih",
}

// ValidateCreate function to validate CreateClusterRequest
func ValidateCreate(req *CreateClusterRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateClusterRequest
func ValidateUpdate(req *UpdateClusterRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
