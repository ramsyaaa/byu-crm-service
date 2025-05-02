package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/communication/response"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type communicationRepository struct {
	db *gorm.DB
}

func NewCommunicationRepository(db *gorm.DB) CommunicationRepository {
	return &communicationRepository{db: db}
}

func (r *communicationRepository) GetAllCommunications(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	accountID int,
) ([]models.Communication, int64, error) {
	var communications []models.Communication
	var total int64

	query := r.db.Model(&models.Communication{}).Where("communications.account_id = ?", accountID).
		Preload("MainCommunication").
		Preload("NextCommunication").
		Preload("Account").
		Preload("Contact").
		Preload("Opportunity")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("communications.note LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("communications.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("communications.created_at <= ?", endDate)
	}

	// Count total sebelum pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&communications).Error
	return communications, total, err
}

func (r *communicationRepository) FindByCommunicationID(id uint) (*models.Communication, error) {
	var communication models.Communication

	query := r.db.
		Model(&models.Communication{}).
		Preload("MainCommunication").
		Preload("NextCommunication").
		Preload("Account").
		Preload("Contact").
		Preload("Opportunity").
		Where("communications.id = ?", id)

	err := query.First(&communication).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &communication, nil
}

func (r *communicationRepository) CreateCommunication(requestBody map[string]string) (*models.Communication, error) {
	communication := models.Communication{
		CommunicationType:   func(s string) *string { return &s }(requestBody["communication_type"]),
		Note:                func(s string) *string { return &s }(requestBody["note"]),
		Date:                func() *time.Time { now := time.Now(); return &now }(),
		StatusCommunication: func(s string) *string { return &s }(requestBody["status_communication"]),
		OpportunityID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["opportunity_id"]),
		AccountID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["account_id"]),
		ContactID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["contact_id"]),
		CreatedBy: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["created_by"]),
	}

	if err := r.db.Create(&communication).Error; err != nil {
		return nil, err
	}

	var newCommunication, _ = r.FindByCommunicationID(communication.ID)

	return newCommunication, nil
}

func (r *communicationRepository) UpdateAccount(requestBody map[string]string, accountID int, userID int) ([]models.Account, error) {
	var account models.Account

	// Cek dulu apakah account dengan ID itu ada
	if err := r.db.First(&account, accountID).Error; err != nil {
		return nil, err
	}

	// Mapping ulang semua field kayak CreateAccount
	updatedAccount := models.Account{
		AccountName:     func(s string) *string { return &s }(requestBody["account_name"]),
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

	// Update semua kolom
	if err := r.db.Model(&account).Updates(updatedAccount).Error; err != nil {
		return nil, err
	}

	// Ambil hasil yang sudah diupdate
	var updatedAccounts []models.Account
	if err := r.db.Where("id = ?", accountID).Find(&updatedAccounts).Error; err != nil {
		return nil, err
	}

	return updatedAccounts, nil
}

func (r *communicationRepository) GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error) {
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

func (r *communicationRepository) FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.AccountResponse, error) {
	var account response.AccountResponse

	query := r.db.
		Model(&models.Account{}).
		Select(`
			accounts.*,
			cities.name AS city_name,
			clusters.name AS cluster_name,
			branches.name AS branch_name,
			regions.name AS region_name,
			areas.name AS area_name
		`).
		Joins("LEFT JOIN cities ON accounts.city = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id").
		Preload("SocialMedias", "subject_type = ?", "App\\Models\\Account").
		Preload("AccountTypeCampusDetail").
		Preload("AccountTypeSchoolDetail").
		Preload("AccountTypeCommunityDetail").
		Preload("AccountCity.Cluster.Branch.Region.Area").
		Preload("AccountFaculties.Faculty").
		Preload("Contacts").
		Preload("AccountMembers", "subject_type = ? AND subject_id = ?", "App\\Models\\Account", id).
		Preload("AccountLectures", "subject_type = ? AND subject_id = ?", "App\\Models\\AccountLecture", id).
		Where("accounts.id = ?", id)

	err := query.First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// --- Validasi akses berdasarkan wilayah dan role ---
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

		var multiTerritories []models.MultipleTerritory
		r.db.Where("user_id = ?", userID).Find(&multiTerritories)

		for _, mt := range multiTerritories {
			var ids []string
			err := json.Unmarshal([]byte(mt.SubjectIDs), &ids)
			if err != nil || len(ids) == 0 {
				continue
			}

			// Helper function to check if ID exists in the slice
			contains := func(slice []string, item string) bool {
				for _, v := range slice {
					if v == item {
						return true
					}
				}
				return false
			}

			switch mt.SubjectType {
			case "App\\Models\\Area":
				if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil &&
					account.AccountCity.Cluster.Branch.Region != nil && account.AccountCity.Cluster.Branch.Region.Area != nil {
					hasAccess = contains(ids, fmt.Sprintf("%d", account.AccountCity.Cluster.Branch.Region.Area.ID))
				}
			case "App\\Models\\Region":
				if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil &&
					account.AccountCity.Cluster.Branch.Region != nil {
					hasAccess = contains(ids, fmt.Sprintf("%d", account.AccountCity.Cluster.Branch.Region.ID))
				}
			case "App\\Models\\Branch":
				if account.AccountCity != nil && account.AccountCity.Cluster != nil && account.AccountCity.Cluster.Branch != nil {
					hasAccess = contains(ids, fmt.Sprintf("%d", account.AccountCity.Cluster.Branch.ID))
				}
			case "App\\Models\\Cluster":
				if account.AccountCity != nil && account.AccountCity.Cluster != nil {
					hasAccess = contains(ids, fmt.Sprintf("%d", account.AccountCity.Cluster.ID))
				}
			}
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

	// --- Kosongkan field berdasarkan kategori akun ---
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

func (r *communicationRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *communicationRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).Updates(fields).Error
}

func (r *communicationRepository) GetAccountVisitCounts(
	filters map[string]string,
	userRole string,
	territoryID int,
	userID int,
) (int64, int64, int64, error) {
	var visitedCount int64
	var notVisitedCount int64
	var totalCount int64

	// Base query
	baseQuery := r.db.Model(&models.Account{}).
		Joins("LEFT JOIN cities ON accounts.City = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			baseQuery = baseQuery.Where(
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
		// Filter utama berdasarkan user role
		switch userRole {
		case "Area":
			baseQuery = baseQuery.Where("areas.id = ?", territoryID)
		case "Regional":
			baseQuery = baseQuery.Where("regions.id = ?", territoryID)
		case "Branch", "Buddies", "DS", "Organic", "YAE":
			baseQuery = baseQuery.Where("branches.id = ?", territoryID)
		case "Admin-Tap", "Cluster":
			baseQuery = baseQuery.Where("clusters.id = ?", territoryID)
		}

		// Tambahan dari multiple_territories
		var multiTerritories []models.MultipleTerritory
		r.db.Where("user_id = ?", userID).Find(&multiTerritories)

		orQuery := r.db.Session(&gorm.Session{}).Model(&models.Account{})

		for _, mt := range multiTerritories {
			var ids []string
			err := json.Unmarshal([]byte(mt.SubjectIDs), &ids)
			if err != nil || len(ids) == 0 {
				continue
			}

			switch mt.SubjectType {
			case "App\\Models\\Area":
				orQuery = orQuery.Or("areas.id IN ?", ids)
			case "App\\Models\\Region":
				orQuery = orQuery.Or("regions.id IN ?", ids)
			case "App\\Models\\Branch":
				orQuery = orQuery.Or("branches.id IN ?", ids)
			case "App\\Models\\Cluster":
				orQuery = orQuery.Or("clusters.id IN ?", ids)
			}
		}

		baseQuery = baseQuery.Or(orQuery)
	}

	if userID > 0 {
		if userRole == "Buddies" || userRole == "DS" {
			baseQuery = baseQuery.Where("accounts.pic = ? OR accounts.pic IS NULL", userID)
		}
	}

	// Hitung total account terlebih dahulu
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return 0, 0, 0, err
	}

	// Ambil data kunjungan
	now := time.Now()
	var visitedAccountIDs []int
	err := r.db.Model(&models.AbsenceUser{}).
		Select("subject_id").
		Where("user_id = ? AND type = ? AND MONTH(clock_in) = ? AND YEAR(clock_in) = ?", userID, "Visit Account", now.Month(), now.Year()).
		Where("subject_id IS NOT NULL").
		Pluck("subject_id", &visitedAccountIDs).Error
	if err != nil {
		return 0, 0, 0, err
	}

	if len(visitedAccountIDs) == 0 {
		// Jika tidak ada akun yang dikunjungi, visited = 0, not visited = total
		return 0, totalCount, totalCount, nil
	}

	// Hitung visited
	visitedQuery := baseQuery.Session(&gorm.Session{}).Where("accounts.id IN ?", visitedAccountIDs)
	if err := visitedQuery.Count(&visitedCount).Error; err != nil {
		return 0, 0, 0, err
	}

	// Sisanya dianggap tidak dikunjungi
	notVisitedCount = totalCount - visitedCount

	return visitedCount, notVisitedCount, totalCount, nil
}
