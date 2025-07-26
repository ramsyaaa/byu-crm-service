package validation

import (
	"byu-crm-service/helper"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	ProgramName  string `json:"program_name" validate:"required"`  // input
	ContractDate string `json:"contract_date" validate:"required"` // input
	AccountID    string `json:"account_id" validate:"required"`    // select

	FirstPartyName        string `json:"first_party_name" validate:"required"`         // input
	FirstPartyPosition    string `json:"first_party_position" validate:"required"`     // input
	FirstPartyPhoneNumber string `json:"first_party_phone_number" validate:"required"` // input
	FirstPartyAddress     string `json:"first_party_address" validate:"required"`      // textarea

	SecondPartyCompany     string `json:"second_party_company" validate:"required"`      // input
	SecondPartyName        string `json:"second_party_name" validate:"required"`         // input
	SecondPartyPosition    string `json:"second_party_position" validate:"required"`     // input
	SecondPartyPhoneNumber string `json:"second_party_phone_number" validate:"required"` // input
	SecondPartyAddress     string `json:"second_party_address" validate:"required"`      // textarea

	Description string `json:"description" validate:"required"` // textarea

	AdditionalSignTitle    *[]string `json:"additional_sign_title"`
	AdditionalSignName     *[]string `json:"additional_sign_name"`
	AdditionalSignPosition *[]string `json:"additional_sign_position"`
}

var validationMessages = map[string]string{
	"program_name.required":  "Nama program wajib diisi",
	"contract_date.required": "Tanggal kontrak wajib diisi",
	"account_id.required":    "Account wajib diisi",

	"first_party_name.required":         "Nama pihak pertama wajib diisi",
	"first_party_position.required":     "Jabatan pihak pertama wajib diisi",
	"first_party_phone_number.required": "Nomor telepon pihak pertama wajib diisi",
	"first_party_address.required":      "Alamat pihak pertama wajib diisi",

	"second_party_company.required":      "Perusahaan pihak kedua wajib diisi",
	"second_party_name.required":         "Nama pihak kedua wajib diisi",
	"second_party_position.required":     "Jabatan pihak kedua wajib diisi",
	"second_party_phone_number.required": "Nomor telepon pihak kedua wajib diisi",
	"second_party_address.required":      "Alamat pihak kedua wajib diisi",

	"description.required": "Deskripsi wajib diisi",
}

var validate = validator.New()

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateAdditional(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)
	if _, err := time.Parse("2006-01-02", req.ContractDate); err != nil {
		errors["contract_date"] = "Format tanggal tidak valid (gunakan YYYY-MM-DD)."
	}

	// Validasi AdditionalSign*
	counts := []int{}
	if req.AdditionalSignTitle != nil {
		counts = append(counts, len(*req.AdditionalSignTitle))
	}
	if req.AdditionalSignName != nil {
		counts = append(counts, len(*req.AdditionalSignName))
	}
	if req.AdditionalSignPosition != nil {
		counts = append(counts, len(*req.AdditionalSignPosition))
	}

	if len(counts) > 0 {
		if req.AdditionalSignTitle != nil || req.AdditionalSignName != nil || req.AdditionalSignPosition != nil {
			// Pastikan semua field tidak nil
			if req.AdditionalSignTitle == nil {
				errors["additional_sign_title"] = "Field title tanda tangan wajib diisi jika field tambahan lainnya diisi."
			}
			if req.AdditionalSignName == nil {
				errors["additional_sign_name"] = "Field nama penanda tangan wajib diisi jika field tambahan lainnya diisi."
			}
			if req.AdditionalSignPosition == nil {
				errors["additional_sign_position"] = "Field jabatan penanda tangan wajib diisi jika field tambahan lainnya diisi."
			}

			// Lanjutkan validasi hanya jika semuanya tidak nil
			if req.AdditionalSignTitle != nil && req.AdditionalSignName != nil && req.AdditionalSignPosition != nil {
				titles := *req.AdditionalSignTitle
				names := *req.AdditionalSignName
				positions := *req.AdditionalSignPosition

				titleCount := len(titles)
				nameCount := len(names)
				posCount := len(positions)

				if titleCount != nameCount {
					errors["additional_sign_title"] = "Jumlah data title tanda tangan dan nama penanda tangan harus sama."
					errors["additional_sign_name"] = "Jumlah data title tanda tangan dan nama penanda tangan harus sama."
				}
				if nameCount != posCount {
					errors["additional_sign_name"] = "Jumlah data nama penanda tangan dan jabatan penanda tangan harus sama."
					errors["additional_sign_position"] = "Jumlah data nama penanda tangan dan jabatan penanda tangan harus sama."
				}
				if titleCount > 3 {
					errors["additional_sign_title"] = "Jumlah tanda tangan tambahan maksimal 3."
				}

				// Cek jika ada data kosong ("")
				for i := 0; i < titleCount && i < nameCount && i < posCount; i++ {
					if strings.TrimSpace(titles[i]) == "" {
						errors["additional_sign_title"] = fmt.Sprintf("Data ke-%d pada title tanda tangan tidak boleh kosong.", i+1)
					}
					if strings.TrimSpace(names[i]) == "" {
						errors["additional_sign_name"] = fmt.Sprintf("Data ke-%d pada nama penanda tangan tidak boleh kosong.", i+1)
					}
					if strings.TrimSpace(positions[i]) == "" {
						errors["additional_sign_position"] = fmt.Sprintf("Data ke-%d pada jabatan penanda tangan tidak boleh kosong.", i+1)
					}
				}
			}
		}

	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
