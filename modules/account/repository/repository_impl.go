package repository

import (
	"byu-crm-service/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Model(&models.Account{}).
		Joins("LEFT JOIN cities ON accounts.Accountcity = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("accounts.account_name LIKE ?", "%"+token+"%").
					Or("accounts.account_code LIKE ?", "%"+token+"%").
					Or("accounts.account_type LIKE ?", "%"+token+"%").
					Or("accounts.account_category LIKE ?", "%"+token+"%").
					Or("cities.name LIKE ?", "%"+token+"%").
					Or("clusters.name LIKE ?", "%"+token+"%").
					Or("branches.name LIKE ?", "%"+token+"%").
					Or("regions.name LIKE ?", "%"+token+"%").
					Or("areas.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply user role and territory filtering
	if userRole != "Super-Admin" && userRole != "HQ" {
		switch userRole {
		case "Area":
			query = query.Where("areas.id = ?", territoryID)
		case "Regional":
			query = query.Where("regions.id = ?", territoryID)
		case "Branch", "Buddies", "DS", "Organic":
			query = query.Where("branches.id = ?", territoryID)
		case "Admin-Tap", "Cluster":
			query = query.Where("clusters.id = ?", territoryID)
		}
	}

	if onlyUserPic && userID > 0 {
		query = query.Where("accounts.pic = ?", userID)
	} else {
		// Additional logic for Buddies / DS role
		if userRole == "Buddies" || userRole == "DS" {
			query = query.Where("accounts.pic = ? OR accounts.pic IS NULL", userID).
				Order(fmt.Sprintf("CASE WHEN accounts.pic = %d THEN 0 ELSE 1 END, accounts.account_name ASC", userID))
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("accounts.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("accounts.created_at <= ?", endDate)
	}

	// Get total count before applying pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if userRole != "Buddies" && userRole != "DS" {
		if orderBy != "" && orderBy != "id" {
			query = query.Order(orderBy + " " + order)
		} else {
			orderClause := "CASE " +
				"WHEN accounts.account_name LIKE '%" + filters["search"] + "%' THEN 1 " +
				"WHEN accounts.account_code LIKE '%" + filters["search"] + "%' THEN 2 " +
				"WHEN cities.name LIKE '%" + filters["search"] + "%' THEN 3 " +
				"WHEN clusters.name LIKE '%" + filters["search"] + "%' THEN 4 " +
				"WHEN branches.name LIKE '%" + filters["search"] + "%' THEN 5 " +
				"WHEN regions.name LIKE '%" + filters["search"] + "%' THEN 6 " +
				"WHEN areas.name LIKE '%" + filters["search"] + "%' THEN 7 " +
				"ELSE 8 END"
			query = query.Order(orderClause)
		}
	}

	// Apply pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&accounts).Error
	return accounts, total, err
}

func (r *accountRepository) CreateAccount(requestBody map[string]string, userID int) ([]models.Account, error) {
	account := models.Account{
		AccountName:     func(s string) *string { return &s }(requestBody["account_name"]),
		AccountImage:    func(s string) *string { return &s }(requestBody["account_image"]),
		AccountType:     func(s string) *string { return &s }(requestBody["account_type"]),
		AccountCategory: func(s string) *string { return &s }(requestBody["account_category"]),
		AccountCode:     func(s string) *string { return &s }(requestBody["account_code"]),
		City: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["city"]),
		ContactName:             func(s string) *string { return &s }(requestBody["contact_name"]),
		EmailAccount:            func(s string) *string { return &s }(requestBody["email_account"]),
		WebsiteAccount:          func(s string) *string { return &s }(requestBody["website_account"]),
		SystemInformasiAkademik: func(s string) *string { return &s }(requestBody["system_informasi_akademik"]),
		Latitude:                func(s string) *string { return &s }(requestBody["latitude"]),
		Longitude:               func(s string) *string { return &s }(requestBody["longitude"]),
		Ownership:               func(s string) *string { return &s }(requestBody["ownership"]),
		Pic:                     func(s string) *string { return &s }(requestBody["pic"]),
		PicInternal:             func(s string) *string { return &s }(requestBody["pic_internal"]),
	}

	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}

	var accounts []models.Account
	if err := r.db.Where("id = ?", account.ID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *accountRepository) UpdateAccount(requestBody map[string]string, accountID int, userID int) ([]models.Account, error) {
	var account models.Account

	// Cek apakah akun dengan accountID ada
	if err := r.db.First(&account, accountID).Error; err != nil {
		return nil, err // Akun tidak ditemukan
	}

	// Menyiapkan data yang akan diupdate
	updateData := map[string]interface{}{}

	if v, exists := requestBody["account_name"]; exists {
		updateData["account_name"] = v
	}
	if v, exists := requestBody["account_image"]; exists {
		updateData["account_image"] = v
	}
	if v, exists := requestBody["account_type"]; exists {
		updateData["account_type"] = v
	}
	if v, exists := requestBody["account_category"]; exists {
		updateData["account_category"] = v
	}
	if v, exists := requestBody["account_code"]; exists {
		updateData["account_code"] = v
	}
	if v, exists := requestBody["city"]; exists {
		updateData["city"] = v
	}
	if v, exists := requestBody["contact_name"]; exists {
		updateData["contact_name"] = v
	}
	if v, exists := requestBody["email_account"]; exists {
		updateData["email_account"] = v
	}
	if v, exists := requestBody["website_account"]; exists {
		updateData["website_account"] = v
	}
	if v, exists := requestBody["system_informasi_akademik"]; exists {
		updateData["system_informasi_akademik"] = v
	}
	if v, exists := requestBody["latitude"]; exists {
		updateData["latitude"] = v
	}
	if v, exists := requestBody["longitude"]; exists {
		updateData["longitude"] = v
	}
	if v, exists := requestBody["ownership"]; exists {
		updateData["ownership"] = v
	}
	if v, exists := requestBody["pic"]; exists {
		updateData["pic"] = v
	}
	if v, exists := requestBody["pic_internal"]; exists {
		updateData["pic_internal"] = v
	}

	// Eksekusi update
	if err := r.db.Model(&account).Where("id = ?", accountID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Mengambil data yang telah diperbarui
	var updatedAccounts []models.Account
	if err := r.db.Where("id = ?", accountID).Find(&updatedAccounts).Error; err != nil {
		return nil, err
	}

	return updatedAccounts, nil
}

func (r *accountRepository) GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error) {
	var accounts []models.Account
	var totalRecords int64

	query := r.db.Model(&models.Account{})

	// Apply filters based on search
	if search != "" {
		query = query.Where("account_name LIKE ?", "%"+search+"%")
	}

	// Apply role-based territorial filters
	switch userRole {
	case "Area":
		query = query.Where("area_id = ?", territoryID)
	case "Regional":
		query = query.Where("region_id = ?", territoryID)
	case "Branch":
		query = query.Where("branch_id = ?", territoryID)
	}

	// Count total records
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, int(totalRecords), nil
}

func (r *accountRepository) FindByAccountName(account_name string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_name = ?", account_name).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*models.Account, error) {
	var account models.Account

	query := r.db.
		Preload("SocialMedias", "subject_type = ?", "App\\Models\\Account").
		Preload("AccountTypeCampusDetail").
		Preload("AccountTypeSchoolDetail").
		Preload("AccountTypeCommunityDetail").
		Preload("AccountCity.Cluster.Branch.Region.Area").
		Preload("AccountFaculties.Faculty")

	err := query.First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Validate based on territory and role
	if userRole != "Super-Admin" && userRole != "HQ" {
		hasAccess := false

		switch userRole {
		case "Area":
			if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil &&
				account.AccountCity.Cluster.Branch.Region != nil && account.AccountCity.Cluster.Branch.Region.Area != nil {
				hasAccess = account.AccountCity.Cluster.Branch.Region.Area.ID == territoryID
			}
		case "Regional":
			if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil &&
				account.AccountCity.Cluster.Branch.Region != nil {
				hasAccess = account.AccountCity.Cluster.Branch.Region.ID == territoryID
			}
		case "Branch", "Buddies", "DS", "Organic", "YAE":
			if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil {
				hasAccess = account.AccountCity.Cluster.Branch.ID == territoryID
			}

		case "Admin-Tap", "Cluster":
			if account.AccountCity != nil && account.AccountCity.Cluster != nil {
				hasAccess = account.AccountCity.Cluster.ID == territoryID
			}
		}

		if !hasAccess {
			return nil, errors.New("unauthorized access to account")
		}

		if userRole == "Buddies" || userRole == "DS" || userRole == "Organic" || userRole == "YAE" {
			if account.Pic == nil || *account.Pic == fmt.Sprintf("%d", userID) {
				hasAccess = true
			} else {
				hasAccess = false
			}
		}

		if !hasAccess {
			return nil, errors.New("unauthorized access to account")
		}
	}

	// Clear Data Based On Category
	if account.AccountCategory != nil {
		switch *account.AccountCategory {
		case "KAMPUS":
			account.AccountTypeSchoolDetail = nil
			account.AccountTypeCommunityDetail = nil
		case "SEKOLAH":
			account.AccountTypeCampusDetail = nil
			account.AccountTypeCommunityDetail = nil
			account.AccountFaculties = nil
		case "KOMUNITAS":
			account.AccountTypeCampusDetail = nil
			account.AccountTypeSchoolDetail = nil
			account.AccountFaculties = nil
		default:
			account.AccountTypeCampusDetail = nil
			account.AccountTypeSchoolDetail = nil
			account.AccountTypeCommunityDetail = nil
			account.AccountFaculties = nil
		}
	} else {
		account.AccountFaculties = nil
	}

	return &account, nil
}

func (r *accountRepository) FindByAccountCode(code string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_code = ?", code).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).Updates(fields).Error
}
