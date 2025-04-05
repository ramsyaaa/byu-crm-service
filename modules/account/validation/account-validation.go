package validation

import (
	"byu-crm-service/modules/account/repository"
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

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Next()
}

// Helper function to validate file extension
func validateFileExtension(file *multipart.FileHeader, allowedExt string) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == "."+allowedExt
}

// Validate mime type for images
func validateMime(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true // Allow empty values
	}

	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
	filename := fl.Field().String()
	ext := strings.ToLower(filepath.Ext(filename))

	// Check if the extension is in the allowed types
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// Validate unique account code
func uniqueAccountCode(fl validator.FieldLevel) bool {
	if accountRepo == nil {
		return true // Skip validation if repository is not set
	}

	code := fl.Field().String()
	account, err := accountRepo.FindByAccountCode(code)
	if err != nil {
		return false // Error occurred during validation
	}

	return account == nil // Return true if account doesn't exist (code is unique)
}
