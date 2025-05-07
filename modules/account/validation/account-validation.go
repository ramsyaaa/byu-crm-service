package validation

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/account/repository"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

type ValidateRequest struct {
	AccountName             string  `json:"account_name" validate:"required"`
	AccountType             string  `json:"account_type" validate:"required"`
	AccountCategory         string  `json:"account_category" validate:"required"`
	AccountCode             *string `json:"account_code"`
	City                    string  `json:"city" validate:"required"`
	ContactName             *string `json:"contact_name"`
	EmailAccount            *string `json:"email_account" validate:"email"`
	WebsiteAccount          *string `json:"website_account"`
	SystemInformasiAkademik *string `json:"system_informasi_akademik"`
	Latitude                *string `json:"latitude"`
	Longitude               *string `json:"longitude"`
	Ownership               *string `json:"ownership"`
	Pic                     *string `json:"pic"`
	PicInternal             *string `json:"pic_internal"`

	ProductAccount []string `json:"product_account"`

	// Social Media
	Category  []string `json:"category" validate:"validate_category_url"`
	Url       []string `json:"url"`
	ContactID []string `json:"contact_id"`

	// account category school
	DiesNatalis             *string `json:"dies_natalis"`
	Extracurricular         *string `json:"extracurricular"`
	FootballFieldBranding   *string `json:"football_field_branding"`
	BasketballFieldBranding *string `json:"basketball_field_branding"`
	WallPaintingBranding    *string `json:"wall_painting_branding"`
	WallMagazineBranding    *string `json:"wall_magazine_branding"`

	// account category campus
	Faculties               []string `json:"faculties"`
	YearLecture             []string `json:"year_lecture"`
	AmountLecture           []string `json:"amount_lecture"`
	Origin                  []string `json:"origin"`
	PercentageOrigin        []string `json:"percentage_origin"`
	OrganizationName        []string `json:"organization_name"`
	PreferenceTechnologies  []string `json:"preference_technologies"`
	MemberNeeds             []string `json:"member_needs"`
	AccessTechnology        *string  `json:"access_technology"`
	Byod                    *string  `json:"byod"`
	ItInfrastructures       []string `json:"it_infrastructures"`
	DigitalCollaborations   []string `json:"digital_collaborations"`
	CampusAdministrationApp *string  `json:"campus_administration_app"`
	ProgramIdentification   []string `json:"program_identification"`
	YearRank                []string `json:"year_rank"`
	Rank                    []string `json:"rank"`
	ProgramStudy            []string `json:"program_study"`

	// account category campus & community
	Year                     []string `json:"year"`
	Amount                   []string `json:"amount"`
	Age                      []string `json:"age"`
	PercentageAge            []string `json:"percentage_age"`
	ScheduleCategory         []string `json:"schedule_category"`
	Title                    []string `json:"title"`
	Date                     []string `json:"date"`
	PotentionalCollaboration *string  `json:"potentional_collaboration"`

	// account category community
	AccountSubtype                  *string  `json:"account_subtype"`
	Group                           *string  `json:"group"`
	GroupName                       *string  `json:"group_name"`
	ProductService                  *string  `json:"product_service"`
	PotentialCollaborationItems     *string  `json:"potential_collaboration_items"`
	Gender                          []string `json:"gender"`
	PercentageGender                []string `json:"percentage_gender"`
	EducationalBackground           []string `json:"educational_background"`
	PercentageEducationalBackground []string `json:"percentage_educational_background"`
	Profession                      []string `json:"profession"`
	PercentageProfession            []string `json:"percentage_profession"`
	Income                          []string `json:"income"`
	PercentageIncome                []string `json:"percentage_income"`
}

var validationMessages = map[string]string{
	"account_name.required":          "Nama akun wajib diisi",
	"account_type.required":          "Tipe akun wajib diisi",
	"account_category.required":      "Kategori akun wajib diisi",
	"city.required":                  "Kota wajib diisi",
	"email_account.email":            "Format email tidak valid",
	"category.validate_category_url": "Kategori dan URL wajib diisi",
}

var validate = validator.New()

func init() {
	validate.RegisterValidation("validate_category_url", validateCategoryAndUrl)
}

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateSchool(req *ValidateRequest, isCreate bool, accountID int, userRole string, territoryID int, userID int) map[string]string {
	errors := make(map[string]string)

	// Validate DiesNatalis
	if req.DiesNatalis != nil && strings.TrimSpace(*req.DiesNatalis) != "" {
		_, err := time.Parse("2006-01-02", *req.DiesNatalis)
		if err != nil {
			errors["dies_natalis"] = "Format Dies Natalis harus berupa tanggal dengan format YYYY-MM-DD"
		}
	}

	// Validate branding fields (kecuali Extracurricular)
	validateBranding := func(fieldValue *string, fieldName string) {
		if fieldValue != nil && strings.TrimSpace(*fieldValue) != "" {
			if *fieldValue != "BELUM BRANDING" && *fieldValue != "SUDAH BRANDING" {
				errors[fieldName] = fmt.Sprintf("%s harus 'BELUM BRANDING' atau 'SUDAH BRANDING'", fieldName)
			}
		}
	}

	validateBranding(req.FootballFieldBranding, "football_field_branding")
	validateBranding(req.BasketballFieldBranding, "basketball_field_branding")
	validateBranding(req.WallPaintingBranding, "wall_painting_branding")
	validateBranding(req.WallMagazineBranding, "wall_magazine_branding")

	// Validate account_code unique
	if req.AccountCode != nil {
		shouldCheckUnique := true

		if !isCreate {
			// Kalau update, ambil data account lama
			existingAccount, err := accountRepo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))
			if err == nil && existingAccount != nil && existingAccount.AccountCode != nil {
				// Bandingkan account_code lama dengan yang baru
				if *existingAccount.AccountCode == *req.AccountCode {
					shouldCheckUnique = false // tidak perlu cek unik kalau sama
				}
			}
		}

		if shouldCheckUnique {
			if ok, msg := accountCodeUnique(*req.AccountCode, true); !ok {
				errors["account_code"] = msg
			}
		}
	}

	// Return nil jika tidak ada error
	if len(errors) == 0 {
		return nil
	}

	return errors
}

func ValidateCommunity(req *ValidateRequest, isCreate bool, accountID int, userRole string, territoryID int, userID int) map[string]string {
	errors := make(map[string]string)

	// Helper untuk validasi jumlah array dan isinya tidak kosong
	validateArrayPair := func(a1, a2 []string, field1, field2 string) {
		if len(a1) != len(a2) {
			errors[field1] = fmt.Sprintf("Jumlah %s dan %s harus sama", field1, field2)
			return
		}
		for i := range a1 {
			if strings.TrimSpace(a1[i]) == "" || strings.TrimSpace(a2[i]) == "" {
				errors[field1] = fmt.Sprintf("%s dan %s tidak boleh kosong", field1, field2)
				return
			}
		}
	}

	// Helper untuk validasi percentage harus total 100%
	validatePercentage := func(percentage []string, field string) {
		total := 0
		for _, p := range percentage {
			if p == "" {
				errors[field] = fmt.Sprintf("Semua data %s harus diisi", field)
				return
			}
			num, err := strconv.Atoi(p)
			if err != nil {
				errors[field] = fmt.Sprintf("Semua data %s harus berupa angka", field)
				return
			}
			total += num
		}
		if total != 100 {
			errors[field] = fmt.Sprintf("Total persentase %s harus 100%%", field)
		}
	}

	// Validate percentage-related fields
	validateArrayPair(req.Age, req.PercentageAge, "age", "percentage_age")
	if len(req.PercentageAge) > 0 {
		validatePercentage(req.PercentageAge, "percentage_age")
	}

	if len(req.Gender) != 2 {
		errors["gender"] = "Gender harus mengirim 2 data yaitu 'Pria' dan 'Wanita'"
	} else {
		// Kalau 2 datanya ada, cek apakah isinya Pria dan Wanita
		validGenders := map[string]bool{
			"Pria":   true,
			"Wanita": true,
		}

		if !validGenders[req.Gender[0]] || !validGenders[req.Gender[1]] {
			errors["gender"] = "Gender harus terdiri dari 'Pria' dan 'Wanita'"
		}
	}

	validateArrayPair(req.EducationalBackground, req.PercentageEducationalBackground, "educational_background", "percentage_educational_background")
	if len(req.PercentageEducationalBackground) > 0 {
		validatePercentage(req.PercentageEducationalBackground, "percentage_educational_background")
	}

	validateArrayPair(req.Profession, req.PercentageProfession, "profession", "percentage_profession")
	if len(req.PercentageProfession) > 0 {
		validatePercentage(req.PercentageProfession, "percentage_profession")
	}

	validateArrayPair(req.Income, req.PercentageIncome, "income", "percentage_income")
	if len(req.PercentageIncome) > 0 {
		validatePercentage(req.PercentageIncome, "percentage_income")
	}

	// Validate ScheduleCategory, Title, Date
	if len(req.ScheduleCategory) != len(req.Title) || len(req.Title) != len(req.Date) {
		errors["schedule_category_title_date"] = "Jumlah Schedule Category, Title, dan Date harus sama"
	} else {
		for i := range req.ScheduleCategory {
			if strings.TrimSpace(req.ScheduleCategory[i]) == "" {
				errors[fmt.Sprintf("schedule_category_%d", i)] = "Schedule Category tidak boleh kosong"
			}
			if strings.TrimSpace(req.Title[i]) == "" {
				errors[fmt.Sprintf("title_%d", i)] = "Title tidak boleh kosong"
			}
			if strings.TrimSpace(req.Date[i]) == "" {
				errors[fmt.Sprintf("date_%d", i)] = "Tanggal tidak boleh kosong"
			} else {
				// Pastikan Date bisa dikonversi ke tanggal format YYYY-MM-DD
				_, err := time.Parse("2006-01-02", req.Date[i])
				if err != nil {
					errors[fmt.Sprintf("date_%d", i)] = "Format tanggal harus YYYY-MM-DD"
				}
			}
		}
	}

	// Validate Account Code jika diisi
	if req.AccountCode != nil && strings.TrimSpace(*req.AccountCode) != "" {
		if isCreate {
			if ok, msg := accountCodeUnique(*req.AccountCode, true); !ok {
				errors["account_code"] = msg
			}
		} else {
			// Update, cek kalau account_code berubah
			oldData, err := accountRepo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				errors["account_code"] = "Gagal mengambil data lama untuk validasi"
			} else if *oldData.AccountCode != *req.AccountCode {
				if ok, msg := accountCodeUnique(*req.AccountCode, true); !ok {
					errors["account_code"] = msg
				}
			}
		}
	}

	// Validate Year dan Amount
	if len(req.Year) > 0 || len(req.Amount) > 0 {
		if len(req.Year) != len(req.Amount) {
			errors["year"] = "Jumlah Year dan Amount harus sama"
		} else {
			for i := range req.Year {
				if strings.TrimSpace(req.Year[i]) == "" || strings.TrimSpace(req.Amount[i]) == "" {
					errors["year"] = "Year dan Amount tidak boleh kosong"
					break
				}
				if _, err := strconv.Atoi(req.Year[i]); err != nil {
					errors["year"] = "Year harus berupa angka"
				}
				if _, err := strconv.Atoi(req.Amount[i]); err != nil {
					errors["amount"] = "Amount harus berupa angka"
				}
			}
		}
	}

	// Return nil jika tidak ada error
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func ValidateCampus(req *ValidateRequest, isCreate bool, accountID int, userRole string, territoryID int, userID int) map[string]string {
	errors := make(map[string]string)

	// Helper untuk validasi jumlah array dan isinya tidak kosong
	validateArrayPair := func(a1, a2 []string, field1, field2 string) {
		if len(a1) != len(a2) {
			errors[field1] = fmt.Sprintf("Jumlah %s dan %s harus sama", field1, field2)
			return
		}
		for i := range a1 {
			if strings.TrimSpace(a1[i]) == "" || strings.TrimSpace(a2[i]) == "" {
				errors[field1] = fmt.Sprintf("%s dan %s tidak boleh kosong", field1, field2)
				return
			}
		}
	}

	// Helper untuk validasi percentage harus total 100%
	validatePercentage := func(percentage []string, field string) {
		total := 0
		for _, p := range percentage {
			if p == "" {
				errors[field] = fmt.Sprintf("Semua data %s harus diisi", field)
				return
			}
			num, err := strconv.Atoi(p)
			if err != nil {
				errors[field] = fmt.Sprintf("Semua data %s harus berupa angka", field)
				return
			}
			total += num
		}
		if total != 100 {
			errors[field] = fmt.Sprintf("Total persentase %s harus 100%%", field)
		}
	}

	// Validate percentage-related fields
	validateArrayPair(req.Age, req.PercentageAge, "age", "percentage_age")
	if len(req.PercentageAge) > 0 {
		validatePercentage(req.PercentageAge, "percentage_age")
	}

	validateArrayPair(req.Rank, req.YearRank, "rank", "year_rank")

	// Validate ScheduleCategory, Title, Date
	if len(req.ScheduleCategory) != len(req.Title) || len(req.Title) != len(req.Date) {
		errors["schedule_category_title_date"] = "Jumlah Schedule Category, Title, dan Date harus sama"
	} else {
		for i := range req.ScheduleCategory {
			if strings.TrimSpace(req.ScheduleCategory[i]) == "" {
				errors[fmt.Sprintf("schedule_category_%d", i)] = "Schedule Category tidak boleh kosong"
			}
			if strings.TrimSpace(req.Title[i]) == "" {
				errors[fmt.Sprintf("title_%d", i)] = "Title tidak boleh kosong"
			}
			if strings.TrimSpace(req.Date[i]) == "" {
				errors[fmt.Sprintf("date_%d", i)] = "Tanggal tidak boleh kosong"
			} else {
				// Pastikan Date bisa dikonversi ke tanggal format YYYY-MM-DD
				_, err := time.Parse("2006-01-02", req.Date[i])
				if err != nil {
					errors[fmt.Sprintf("date_%d", i)] = "Format tanggal harus YYYY-MM-DD"
				}
			}
		}
	}

	// Validate Account Code jika diisi
	if req.AccountCode != nil {
		shouldCheckUnique := true

		if !isCreate {
			// Kalau update, ambil data account lama
			existingAccount, err := accountRepo.FindByAccountID(uint(accountID), userRole, uint(territoryID), uint(userID))
			if err == nil && existingAccount != nil && existingAccount.AccountCode != nil {
				// Bandingkan account_code lama dengan yang baru
				if *existingAccount.AccountCode == *req.AccountCode {
					shouldCheckUnique = false // tidak perlu cek unik kalau sama
				}
			}
		}

		if shouldCheckUnique {
			if ok, msg := accountCodeUnique(*req.AccountCode, true); !ok {
				errors["account_code"] = msg
			}
		}
	}

	// Validate Year dan Amount
	if len(req.Year) > 0 || len(req.Amount) > 0 {
		if len(req.Year) != len(req.Amount) {
			errors["year"] = "Jumlah Year dan Amount harus sama"
		} else {
			for i := range req.Year {
				if strings.TrimSpace(req.Year[i]) == "" || strings.TrimSpace(req.Amount[i]) == "" {
					errors["year"] = "Year dan Amount tidak boleh kosong"
					break
				}
				if _, err := strconv.Atoi(req.Year[i]); err != nil {
					errors["year"] = "Year harus berupa angka"
					break
				}
				if _, err := strconv.Atoi(req.Amount[i]); err != nil {
					errors["amount"] = "Amount harus berupa angka"
					break
				}
			}
		}
	}

	// Validate BYOD
	if strings.TrimSpace(*req.Byod) != "" {
		if *req.Byod != "0" && *req.Byod != "1" {
			errors["byod"] = "BYOD harus berupa '0' atau '1'"
		}
	}

	// Return nil jika tidak ada error
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func accountCodeUnique(code string, required bool) (bool, string) {
	trimmedCode := strings.TrimSpace(code)

	if required && trimmedCode == "" {
		return false, "Kode akun wajib diisi"
	}

	if trimmedCode == "" {
		return true, ""
	}

	if accountRepo == nil {
		return true, ""
	}

	account, err := accountRepo.FindByAccountCode(trimmedCode)
	if err != nil {
		return true, ""
	}

	if account != nil {
		return false, "Kode akun sudah digunakan"
	}

	return true, ""
}

func validateCategoryAndUrl(fl validator.FieldLevel) bool {
	parent := fl.Parent()

	categoryField := parent.FieldByName("Category")
	urlField := parent.FieldByName("Url")

	if !categoryField.IsValid() || !urlField.IsValid() {
		return true
	}

	categories, ok1 := categoryField.Interface().([]string)
	urls, ok2 := urlField.Interface().([]string)

	if !ok1 || !ok2 {
		return false
	}

	if len(categories) != len(urls) {
		return false
	}

	for i := range categories {
		if strings.TrimSpace(categories[i]) == "" || strings.TrimSpace(urls[i]) == "" {
			return false
		}
	}

	return true
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
