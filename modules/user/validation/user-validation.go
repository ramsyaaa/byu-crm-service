package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Biar validasi error pakai nama dari `json` tag, bukan nama field Go
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return name
	})
}

// UpdateProfileRequest struct for creating a faculty
type UpdateProfileRequest struct {
	Name            string `json:"name" validate:"required"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":             "Nama harus diisi",
	"old_password.required":     "Password lama harus diisi",
	"new_password.required":     "Password baru harus diisi",
	"new_password.min":          "Password baru minimal 8 karakter",
	"confirm_password.required": "Konfirmasi password harus diisi",
	"confirm_password.min":      "Konfirmasi password minimal 8 karakter",
}

// ValidateUpdate function to validate UpdateProfileRequest
func ValidateUpdate(req *UpdateProfileRequest) map[string]string {
	errorsMap := make(map[string]string)

	// Validasi basic dari struct (hanya validasi 'name')
	if err := validate.StructPartial(req, "Name"); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			key := fmt.Sprintf("%s.%s", err.Field(), err.Tag())
			if msg, ok := validationMessages[strings.ToLower(key)]; ok {
				errorsMap[err.Field()] = msg
			} else {
				errorsMap[err.Field()] = err.Error()
			}
		}
	}

	// Cek jika salah satu password ada yang diisi
	if req.OldPassword != "" || req.NewPassword != "" || req.ConfirmPassword != "" {
		// Validasi old_password wajib
		if req.OldPassword == "" {
			errorsMap["old_password"] = validationMessages["old_password.required"]
		}

		// Validasi new_password wajib dan minimal 8 karakter
		if req.NewPassword == "" {
			errorsMap["new_password"] = validationMessages["new_password.required"]
		} else if len(req.NewPassword) < 8 {
			errorsMap["new_password"] = validationMessages["new_password.min"]
		}

		// Validasi confirm_password wajib dan minimal 8 karakter
		if req.ConfirmPassword == "" {
			errorsMap["confirm_password"] = validationMessages["confirm_password.required"]
		} else if len(req.ConfirmPassword) < 8 {
			errorsMap["confirm_password"] = validationMessages["confirm_password.min"]
		}

		// Validasi confirm_password harus sama dengan new_password
		if req.NewPassword != "" && req.ConfirmPassword != "" && req.NewPassword != req.ConfirmPassword {
			errorsMap["confirm_password"] = "Konfirmasi password harus sama dengan password baru"
		}
	}

	if len(errorsMap) > 0 {
		return errorsMap
	}

	return nil
}
