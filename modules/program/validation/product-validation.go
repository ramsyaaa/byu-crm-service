package validation

import (
	"byu-crm-service/helper"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	ProductName     string  `json:"product_name" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	ProductCategory string  `json:"product_category" validate:"required"`
	ProductType     string  `json:"product_type" validate:"required"`
	Bid             *string `json:"bid"`
	Price           *string `json:"price"`
	KeyVisual       *string `json:"key_visual"`
	AdditionalFile  *string `json:"additional_file"`
	QuotaValue      *string `json:"quota_value"`
	ValidityValue   *string `json:"validity_value"`
	ValidityUnit    *string `json:"validity_unit"`

	EligibilityCategory []string            `json:"eligibility_category"`
	EligibilityType     []string            `json:"eligibility_type"`
	EligibilityLocation map[string][]string `json:"eligibility_location"`
}

var validationMessages = map[string]string{
	"product_name.required":     "Nama produk wajib diisi",
	"product_type.required":     "Tipe produk wajib diisi",
	"product_category.required": "Kategori produk wajib diisi",
	"description.required":      "Deskripsi wajib diisi",
}

var validate = validator.New()

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateVoucherPerdana(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)

	// Validasi Bid (wajib)
	if req.Bid == nil || strings.TrimSpace(*req.Bid) == "" {
		errors["bid"] = "Bid wajib diisi"
	}

	// Validasi Price (wajib dan harus berupa angka)
	if req.Price == nil || strings.TrimSpace(*req.Price) == "" {
		errors["price"] = "Harga wajib diisi"
	} else {
		cleanPrice := strings.ReplaceAll(*req.Price, "Rp", "")
		cleanPrice = strings.ReplaceAll(cleanPrice, ".", "")
		cleanPrice = strings.TrimSpace(cleanPrice)

		if _, err := strconv.Atoi(cleanPrice); err != nil {
			errors["price"] = "Harga tidak valid, gunakan format seperti 10000 tanpa simbol atau pemisah"
		}
	}

	// Validasi Key Visual (wajib dan harus berupa gambar base64)
	// if req.KeyVisual == nil || strings.TrimSpace(*req.KeyVisual) == "" {
	// 	errors["key_visual"] = "Key Visual wajib diisi"
	// } else {
	// 	data := *req.KeyVisual
	// 	if !strings.HasPrefix(data, "data:image/") || !strings.Contains(data, ";base64,") {
	// 		errors["key_visual"] = "Key Visual harus berupa gambar base64 yang valid"
	// 	} else {
	// 		base64Data := strings.Split(data, ",")[1]
	// 		if _, err := base64.StdEncoding.DecodeString(base64Data); err != nil {
	// 			errors["key_visual"] = "Format base64 Key Visual tidak valid"
	// 		}
	// 	}
	// }

	// Validasi Quota Value (opsional tapi jika diisi harus bisa dikonversi ke float)
	if req.QuotaValue != nil && strings.TrimSpace(*req.QuotaValue) != "" {
		if _, err := strconv.ParseFloat(*req.QuotaValue, 64); err != nil {
			errors["quota_value"] = "Kuota harus berupa angka desimal yang valid"
		}
	}

	// Validasi Validity Value (wajib dan harus bisa dikonversi ke float)
	if req.ValidityValue == nil || strings.TrimSpace(*req.ValidityValue) == "" {
		errors["validity_value"] = "Masa berlaku wajib diisi"
	} else {
		if _, err := strconv.ParseFloat(*req.ValidityValue, 64); err != nil {
			errors["validity_value"] = "Masa berlaku harus berupa angka desimal yang valid"
		}
	}

	// Validasi Validity Unit (wajib)
	if req.ValidityUnit == nil || strings.TrimSpace(*req.ValidityUnit) == "" {
		errors["validity_unit"] = "Satuan masa berlaku wajib diisi"
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func ValidateSolutionLbo(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)

	// if req.AdditionalFile == nil || strings.TrimSpace(*req.AdditionalFile) == "" {
	// 	errors["additional_file"] = "File wajib diisi"
	// } else {
	// 	data := *req.AdditionalFile

	// 	// Jika mengandung prefix data URI (opsional), ambil bagian base64-nya
	// 	if strings.Contains(data, ";base64,") {
	// 		parts := strings.SplitN(data, ",", 2)
	// 		if len(parts) != 2 {
	// 			errors["additional_file"] = "Format base64 tidak valid"
	// 		} else {
	// 			data = parts[1]
	// 		}
	// 	}

	// 	if _, err := base64.StdEncoding.DecodeString(data); err != nil {
	// 		errors["additional_file"] = "File harus berupa base64 yang valid"
	// 	}
	// }

	if len(errors) == 0 {
		return nil
	}
	return errors
}

func ValidateHousehold(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)

	// Validasi Price (wajib dan harus berupa angka)
	if req.Price == nil || strings.TrimSpace(*req.Price) == "" {
		errors["price"] = "Harga wajib diisi"
	} else {
		cleanPrice := strings.ReplaceAll(*req.Price, "Rp", "")
		cleanPrice = strings.ReplaceAll(cleanPrice, ".", "")
		cleanPrice = strings.TrimSpace(cleanPrice)

		if _, err := strconv.Atoi(cleanPrice); err != nil {
			errors["price"] = "Harga tidak valid, gunakan format seperti 10000 tanpa simbol atau pemisah"
		}
	}

	// Validasi Key Visual (wajib dan harus berupa gambar base64)
	// if req.KeyVisual == nil || strings.TrimSpace(*req.KeyVisual) == "" {
	// 	errors["key_visual"] = "Key Visual wajib diisi"
	// } else {
	// 	data := *req.KeyVisual
	// 	if !strings.HasPrefix(data, "data:image/") || !strings.Contains(data, ";base64,") {
	// 		errors["key_visual"] = "Key Visual harus berupa gambar base64 yang valid"
	// 	} else {
	// 		base64Data := strings.Split(data, ",")[1]
	// 		if _, err := base64.StdEncoding.DecodeString(base64Data); err != nil {
	// 			errors["key_visual"] = "Format base64 Key Visual tidak valid"
	// 		}
	// 	}
	// }

	if len(errors) == 0 {
		return nil
	}

	return errors
}
