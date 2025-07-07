package validation

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/user/repository"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var UserRepo repository.UserRepository

func SetUserRepository(repo repository.UserRepository) {
	UserRepo = repo
}

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

type ValidateRequest struct {
	Name            string  `json:"name" validate:"required"`
	Email           string  `json:"email" validate:"required"`
	UserType        string  `json:"user_type" validate:"required"`
	Msisdn          *string `json:"msisdn"`
	OutletIDDigipos *string `json:"outlet_id_digipos"`
	NamiAgentID     *string `json:"nami_agent_id"`
	RoleID          string  `json:"role_id"`
	AreaID          *string `json:"area_id"`
	RegionID        *string `json:"region_id"`
	BranchID        *string `json:"branch_id"`
	ClusterID       *string `json:"cluster_id"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirm_password"`
}

// UpdateProfileRequest struct for creating a faculty
type UpdateProfileRequest struct {
	Name            string `json:"name" validate:"required"`
	Msisdn          string `json:"msisdn" validate:"required"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":             "Nama harus diisi",
	"email.required":            "Email harus diisi",
	"password.required":         "Password harus diisi",
	"password.min":              "Password minimal 8 karakter",
	"old_password.required":     "Password lama harus diisi",
	"new_password.required":     "Password baru harus diisi",
	"new_password.min":          "Password baru minimal 8 karakter",
	"confirm_password.required": "Konfirmasi password harus diisi",
	"confirm_password.min":      "Konfirmasi password minimal 8 karakter",
	"msisdn.required":           "MSISDN harus diisi",
}

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func AdditionalValidate(req *ValidateRequest, userID int) map[string]string {
	errors := make(map[string]string)

	// Validasi RoleID (harus termasuk dalam daftar role yang diperbolehkan)
	allowedRoleIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	valid := false
	for _, id := range allowedRoleIDs {
		if req.RoleID == fmt.Sprintf("%d", id) {
			valid = true
			break
		}
	}
	if !valid {
		errors["role_id"] = "Role tidak valid"
	}

	// Jika role perlu validasi wilayah
	if req.RoleID != "1" && req.RoleID != "2" && req.RoleID != "12" {
		switch req.RoleID {
		case "3":
			if req.AreaID == nil || strings.TrimSpace(*req.AreaID) == "" {
				errors["area_id"] = "Area harus dipilih"
			}
		case "4":
			if req.RegionID == nil || strings.TrimSpace(*req.RegionID) == "" {
				errors["region_id"] = "Region harus dipilih"
			}
		case "5", "7", "9", "10", "11":
			if req.BranchID == nil || strings.TrimSpace(*req.BranchID) == "" {
				errors["branch_id"] = "Branch harus dipilih"
			}
		case "6", "8":
			if req.ClusterID == nil || strings.TrimSpace(*req.ClusterID) == "" {
				errors["cluster_id"] = "Cluster harus dipilih"
			}
		}
	}

	if userID == 0 {
		if req.Password == nil || strings.TrimSpace(*req.Password) == "" || req.ConfirmPassword == nil || strings.TrimSpace(*req.ConfirmPassword) == "" {
			errors["password"] = validationMessages["password.required"]
		}
	}

	// Validasi password hanya jika create atau saat password/confirm_password diisi
	if userID == 0 || req.Password != nil || req.ConfirmPassword != nil {
		if req.Password == nil || strings.TrimSpace(*req.Password) == "" {
			errors["password"] = validationMessages["password.required"]
		} else if len(*req.Password) < 8 {
			errors["password"] = validationMessages["password.min"]
		}

		if req.ConfirmPassword == nil || strings.TrimSpace(*req.ConfirmPassword) == "" {
			errors["confirm_password"] = validationMessages["confirm_password.required"]
		} else if len(*req.ConfirmPassword) < 8 {
			errors["confirm_password"] = validationMessages["confirm_password.min"]
		}

		if req.Password != nil && req.ConfirmPassword != nil && *req.Password != *req.ConfirmPassword {
			errors["confirm_password"] = "Konfirmasi password harus sama dengan password baru"
		}
	}

	// Validasi email
	existingUser, err := UserRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil && existingUser.Email != "" {
		if userID == 0 || int(existingUser.ID) != userID {
			errors["email"] = "Email sudah digunakan"
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// ValidateProfile function to validate UpdateProfileRequest
func ValidateProfile(req *UpdateProfileRequest) map[string]string {
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

	if err := validate.StructPartial(req, "Msisdn"); err != nil {
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
