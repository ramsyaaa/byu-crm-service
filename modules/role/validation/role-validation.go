package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateRoleRequest struct for creating a Area
type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required"`
	PermissionIDs []string `json:"permission_id" validate:"required,min=1,dive,required"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name" validate:"required"`
	PermissionIDs []string `json:"permission_id" validate:"required,min=1,dive,required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":            "Nama role harus diisi",
	"permission_id.required":   "Permission wajib diisi",
	"permission_id.min":        "Minimal harus memiliki 1 permission",
	"permission_id.*.required": "Setiap permission tidak boleh kosong",
}

// ValidateCreate function to validate CreateRoleRequest
func ValidateCreate(req *CreateRoleRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateRoleRequest
func ValidateUpdate(req *UpdateRoleRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
