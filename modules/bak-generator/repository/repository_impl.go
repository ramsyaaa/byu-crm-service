package repository

import (
	"byu-crm-service/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type bakGeneratorRepository struct {
	db *gorm.DB
}

func NewBakGeneratorRepository(db *gorm.DB) BakGeneratorRepository {
	return &bakGeneratorRepository{db: db}
}

func (r *bakGeneratorRepository) CreateBak(requestBody map[string]interface{}, user_id uint) error {

	fmt.Println("account_id", requestBody["account_id"])
	programName, _ := requestBody["program_name"].(string)
	contractDateStr, _ := requestBody["contract_date"].(string)
	var accountID uint
	switch v := requestBody["account_id"].(type) {
	case float64:
		accountID = uint(v)
	case int:
		accountID = uint(v)
	case int64:
		accountID = uint(v)
	case string:
		// Jika dikirim dari form input string
		idParsed, err := strconv.Atoi(v)
		if err == nil {
			accountID = uint(idParsed)
		}
	default:
		accountID = 0 // fallback
	}
	userID := user_id

	// First party
	firstPartyName, _ := requestBody["first_party_name"].(string)
	firstPartyPosition, _ := requestBody["first_party_position"].(string)
	firstPartyPhoneNumber, _ := requestBody["first_party_phone_number"].(string)
	firstPartyAddress, _ := requestBody["first_party_address"].(string)

	// Second party
	secondPartyCompany, _ := requestBody["second_party_company"].(string)
	secondPartyName, _ := requestBody["second_party_name"].(string)
	secondPartyPosition, _ := requestBody["second_party_position"].(string)
	secondPartyPhoneNumber, _ := requestBody["second_party_phone_number"].(string)
	secondPartyAddress, _ := requestBody["second_party_address"].(string)

	// Description
	description, _ := requestBody["description"].(string)

	// Additional signs (nullable)
	var additionalSignTitle *string
	if val, exists := requestBody["additional_sign_title"]; exists && val != nil {
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			str := string(jsonBytes)
			additionalSignTitle = &str
		}
	}

	var additionalSignName *string
	if val, exists := requestBody["additional_sign_name"]; exists && val != nil {
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			str := string(jsonBytes)
			additionalSignName = &str
		}
	}

	var additionalSignPosition *string
	if val, exists := requestBody["additional_sign_position"]; exists && val != nil {
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			str := string(jsonBytes)
			additionalSignPosition = &str
		}
	}

	// Parse date if necessary
	contractDate, _ := time.Parse("2006-01-02", contractDateStr)

	account := models.BakFile{
		ProgramName:            programName,
		ContractDate:           contractDate,
		AccountID:              accountID,
		UserID:                 userID,
		FirstPartyName:         firstPartyName,
		FirstPartyPosition:     firstPartyPosition,
		FirstPartyPhoneNumber:  firstPartyPhoneNumber,
		FirstPartyAddress:      firstPartyAddress,
		SecondPartyCompany:     secondPartyCompany,
		SecondPartyName:        secondPartyName,
		SecondPartyPosition:    secondPartyPosition,
		SecondPartyPhoneNumber: secondPartyPhoneNumber,
		SecondPartyAddress:     secondPartyAddress,
		Description:            description,
		AdditionalSignTitle:    additionalSignTitle,
		AdditionalSignName:     additionalSignName,
		AdditionalSignPosition: additionalSignPosition,
	}

	if err := r.db.Create(&account).Error; err != nil {
		return err
	}
	return nil
}

func (r *bakGeneratorRepository) GetAllBak(limit int, paginate bool, page int, filters map[string]string) ([]models.BakFile, int, error) {
	var baks []models.BakFile
	var total int64 // gunakan int64 untuk hasil Count()

	db := r.db.Preload("Account")

	// Filter search
	if search := filters["search"]; search != "" {
		searchPattern := "%" + search + "%"
		db = db.Where(
			r.db.Where("program_name LIKE ?", searchPattern).
				Or("first_party_name LIKE ?", searchPattern).
				Or("second_party_name LIKE ?", searchPattern),
		)
	}

	// Filter contract_date range
	if start := filters["start_date"]; start != "" {
		db = db.Where("contract_date >= ?", start)
	}
	if end := filters["end_date"]; end != "" {
		db = db.Where("contract_date <= ?", end)
	}

	// Order by
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy != "" && order != "" {
		db = db.Order(orderBy + " " + order)
	}

	// Hitung total sebelum limit & offset
	if err := db.Model(&models.BakFile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan paginasi / limit
	if paginate {
		if limit > 0 {
			offset := (page - 1) * limit
			db = db.Limit(limit).Offset(offset)
		}
	} else {
		if limit > 0 {
			db = db.Limit(limit)
		}
		// limit == 0 → ambil semua
	}

	// Eksekusi query data
	if err := db.Find(&baks).Error; err != nil {
		return nil, 0, err
	}

	return baks, int(total), nil
}

func (r *bakGeneratorRepository) GetBakByID(id uint) (*models.BakFile, error) {
	var bak models.BakFile
	if err := r.db.Preload("Account").First(&bak, id).Error; err != nil {
		return nil, err
	}
	return &bak, nil
}
