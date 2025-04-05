package validation

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var accountRepo repository.AccountRepository

func SetAccountRepository(repo repository.AccountRepository) {
	accountRepo = repo
}

type UploadRequest struct {
	UserID string `form:"user_id" validate:"required"`
}

type CreateRequest struct {
	AccountImage            *string `json:"account_image" validate:"omitempty,mime"` // Custom mime validation
	AccountName             *string `json:"account_name" validate:"required"`
	AccountType             *string `json:"account_type" validate:"required"`
	AccountCategory         *string `json:"account_category" validate:"required"`
	AccountCode             *string `json:"account_code" validate:"required,unique_account_code"`
	City                    *string `json:"city" validate:"required"`
	ContactName             *string `json:"contact_name"`
	EmailAccount            *string `json:"email_account" validate:"email"`
	WebsiteAccount          *string `json:"website_account"`
	SystemInformasiAkademik *string `json:"system_informasi_akademik"`
	Latitude                *string `json:"latitude"`
	Longitude               *string `json:"longitude"`
	Ownership               *string `json:"ownership"`
	Pic                     *string `json:"pic"`
	PicInternal             *string `json:"pic_internal"`
}

var validationMessages = map[string]string{
	"AccountImage.mime":               "Format gambar tidak valid, yang diizinkan: jpg, jpeg, png, gif",
	"AccountName.required":            "Nama akun wajib diisi",
	"AccountType.required":            "Tipe akun wajib diisi",
	"AccountCategory.required":        "Kategori akun wajib diisi",
	"AccountCode.required":            "Kode akun wajib diisi",
	"AccountCode.unique_account_code": "Kode akun harus unik",
	"City.required":                   "Kota wajib diisi",
	"EmailAccount.email":              "Format email tidak valid",
	"Ownership.required":              "Ownership wajib diisi",
	"Pic.required":                    "PIC wajib diisi",
	"PicInternal.required":            "PIC internal wajib diisi",
}

var validate = validator.New()

func init() {
	validate.RegisterValidation("mime", validateMime)
	validate.RegisterValidation("unique_account_code", uniqueAccountCode)
}

func ValidateAccountRequest(c *fiber.Ctx) error {
	// Check if file exists in the request
	file, err := c.FormFile("file_csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate file extension
	if !validateFileExtension(file, "csv") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only CSV files are allowed",
		})
	}

	// Validate user_id
	var request UploadRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

<<<<<<< HEAD
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
=======
	localValidate := validator.New()
	localValidate.RegisterValidation("file_extension", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "csv"
	})

	if err := localValidate.Struct(request); err != nil {
>>>>>>> 0320bde8277687a5ef50585037e6f48cc6e121b7
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Next()
}

<<<<<<< HEAD
// Helper function to validate file extension
func validateFileExtension(file *multipart.FileHeader, allowedExt string) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == "."+allowedExt
=======
func ValidateCreate(req *CreateRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}

func validateMime(fl validator.FieldLevel) bool {
	// Get the file extension from the field value (assuming it's a string representing the file name)
	fileName := fl.Field().String()

	// Allowed file extensions
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}

	// Check if the file extension is in the allowed list
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

func uniqueAccountCode(fl validator.FieldLevel) bool {
	accountCode := fl.Field().String()

	// Use the FindByAccountCode method to check if the account code already exists
	account, err := accountRepo.FindByAccountCode(accountCode)
	if err != nil && err.Error() != "record not found" {
		// If there's an error other than "record not found", validation fails
		return false
	}

	// If the account is found, it means the account code is not unique
	return account == nil
>>>>>>> 0320bde8277687a5ef50585037e6f48cc6e121b7
}
