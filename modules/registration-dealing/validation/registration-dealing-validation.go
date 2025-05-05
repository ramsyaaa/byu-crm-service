package validation

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/registration-dealing/repository"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var registrationDealingRepo repository.RegistrationDealingRepository

func SetRegistrationDealingRepository(repo repository.RegistrationDealingRepository) {
	registrationDealingRepo = repo
}

type ValidateRequest struct {
	PhoneNumber          string  `json:"phone_number" validate:"required"`
	AccountID            string  `json:"account_id" validate:"required"`
	CustomerName         string  `json:"customer_name" validate:"required"`
	EventName            string  `json:"event_name" validate:"required"`
	WhatsappNumber       string  `json:"whatsapp_number" validate:"required"`
	Class                string  `json:"class" validate:"required"`
	Email                string  `json:"email" validate:"required"`
	SchoolType           string  `json:"school_type" validate:"required"`
	RegistrationEvidence *string `json:"registration_evidence"`
}

var validationMessages = map[string]string{
	"phone_number.required":    "Nomor telepon harus diisi",
	"account_id.required":      "Account ID harus diisi",
	"customer_name.required":   "Nama pelanggan harus diisi",
	"event_name.required":      "Nama event harus diisi",
	"whatsapp_number.required": "Nomor WhatsApp harus diisi",
	"class.required":           "Kelas harus diisi",
	"email.required":           "Email harus diisi",
	"school_type.required":     "Tipe sekolah harus diisi",
}

var validate = validator.New()

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidatePhoneNumber(phone_number string) (bool, map[string]string) {
	trimmedPhoneNumber := strings.TrimSpace(phone_number)
	fmt.Println("Trimmed Phone Number:", trimmedPhoneNumber)

	if trimmedPhoneNumber == "" {
		fmt.Println("Phone number is empty")
		return false, map[string]string{"phone_number": "Nomor telepon harus diisi"}
	}

	if registrationDealingRepo == nil {
		fmt.Println("RegistrationDealingRepository is not set")
		return false, nil
	}

	registrationDealing, err := registrationDealingRepo.FindByPhoneNumber(trimmedPhoneNumber)
	if err != nil {
		fmt.Println("Error checking phone number:", err)
		return true, nil
	}

	fmt.Println(registrationDealing)
	if registrationDealing != nil {
		return false, map[string]string{"phone_number": "Nomor ponsel sudah digunakan"}
	}

	return true, nil
}
