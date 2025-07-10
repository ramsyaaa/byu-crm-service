package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/response"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetAllAccounts(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	userRole string,
	territoryID int,
	userID int,
	onlyUserPic bool,
	excludeVisited bool,
) ([]response.AccountResponse, int64, error) {
	var accounts []response.AccountResponse
	var total int64

	query := r.db.Model(&models.Account{}).
		Select(`
			accounts.*,
			cities.name AS city_name,
			clusters.name AS cluster_name,
			branches.name AS branch_name,
			regions.name AS region_name,
			areas.name AS area_name
		`).
		Joins("LEFT JOIN cities ON accounts.City = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
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

	// Apply accountCategory filter if exists and not empty
	if accountCategory, exists := filters["account_category"]; exists && accountCategory != "" {
		query = query.Where("accounts.account_category = ?", accountCategory)
	}

	// Apply accountType filter if exists and not empty
	if accountType, exists := filters["account_type"]; exists && accountType != "" {
		query = query.Where("accounts.account_type = ?", accountType)
	}

	// Apply user role and territory filtering
	if userRole != "Super-Admin" && userRole != "HQ" {
		// Filter utama berdasarkan user role
		switch userRole {
		case "Area":
			query = query.Where("areas.id = ?", territoryID)
		case "Regional":
			query = query.Where("regions.id = ?", territoryID)
		case "Branch", "Buddies", "DS", "Organic", "YAE":
			fmt.Println("branch id", territoryID)
			query = query.Where("branches.id = ?", territoryID)
		case "Admin-Tap", "Cluster":
			query = query.Where("clusters.id = ?", territoryID)
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

		query = query.Or(orQuery)
	}

	// Only User PIC
	if onlyUserPic && userID > 0 {
		query = query.Where("accounts.pic = ?", userID)
	} else {
		if userRole == "Buddies" || userRole == "DS" || userRole == "YAE" {
			query = query.
				Where("accounts.pic = ? OR accounts.pic IS NULL OR accounts.pic = ''", userID).
				Order(fmt.Sprintf("CASE WHEN accounts.pic = '%d' THEN 0 ELSE 1 END, accounts.account_name ASC", userID))
		}
	}

	// Exclude visited accounts
	if excludeVisited && userID > 0 {
		now := time.Now()
		var visitedAccountIDs []int

		r.db.Model(&models.AbsenceUser{}).
			Select("subject_id").
			Where("user_id = ? AND type = ? AND MONTH(clock_in) = ? AND YEAR(clock_in) = ?", userID, "Visit Account", now.Month(), now.Year()).
			Where("subject_id IS NOT NULL").
			Pluck("subject_id", &visitedAccountIDs)

		if len(visitedAccountIDs) > 0 {
			query = query.Where("accounts.id NOT IN ?", visitedAccountIDs)
		}
	}

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("accounts.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("accounts.created_at <= ?", endDate)
	}

	if onlySkulID, exists := filters["only_skulid"]; exists && onlySkulID == "1" {
		query = query.Where("accounts.is_skulid = ?", 1)
	}

	if isPriority, exists := filters["is_priority"]; exists && isPriority == "1" {

		if priorityStr, exists := filters["priority"]; exists && priorityStr != "" {
			priorities := strings.Split(priorityStr, ",")
			query = query.Where("accounts.priority IN ?", priorities)
		}

		// if userRole == "Buddies" || userRole == "DS" || userRole == "Organic" || userRole == "YAE" {
		// 	if territoryID == 84 {
		// 		var cityIDs []int
		// 		// Ambil semua ID kota yang termasuk dalam cluster_id = 243
		// 		err := r.db.Model(&models.City{}).
		// 			Where("cluster_id = ?", 243).
		// 			Pluck("id", &cityIDs).Error

		// 		if err == nil && len(cityIDs) > 0 {
		// 			// Tambahkan pengecualian OR untuk city id tersebut
		// 			query = query.Or("accounts.City IN ?", cityIDs)
		// 		}
		// 	}
		// }
	}

	// Count total before pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if userRole != "Buddies" && userRole != "DS" {
		if orderBy != "" {
			query = query.Order(orderBy + " " + order)
		} else if filters["search"] != "" {
			search := filters["search"]
			orderClause := "CASE " +
				"WHEN accounts.account_name LIKE '%" + search + "%' THEN 1 " +
				"WHEN accounts.account_code LIKE '%" + search + "%' THEN 2 " +
				"WHEN cities.name LIKE '%" + search + "%' THEN 3 " +
				"WHEN clusters.name LIKE '%" + search + "%' THEN 4 " +
				"WHEN branches.name LIKE '%" + search + "%' THEN 5 " +
				"WHEN regions.name LIKE '%" + search + "%' THEN 6 " +
				"WHEN areas.name LIKE '%" + search + "%' THEN 7 " +
				"ELSE 8 END"
			query = query.Order(orderClause)
		}
	}

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&accounts).Error
	if err != nil {
		return accounts, total, err
	}

	// Cek jika latitude dan longitude diberikan
	latStr, latOk := filters["latitude"]
	lonStr, lonOk := filters["longitude"]

	if latOk && lonOk && latStr != "" && lonStr != "" {
		userLat, err1 := strconv.ParseFloat(latStr, 64)
		userLon, err2 := strconv.ParseFloat(lonStr, 64)

		if err1 == nil && err2 == nil {
			for i, acc := range accounts {
				newDistance := "-"
				accounts[i].Distance = &newDistance

				if acc.Latitude != nil && acc.Longitude != nil {
					lat, err1 := strconv.ParseFloat(*acc.Latitude, 64)
					lon, err2 := strconv.ParseFloat(*acc.Longitude, 64)
					if err1 == nil && err2 == nil {
						distance := haversine(userLat, userLon, lat, lon)

						var distanceStr string
						if distance < 1000 {
							distanceStr = fmt.Sprintf("%.2f m", distance)
						} else {
							distanceStr = fmt.Sprintf("%.2f km", distance/1000)
						}

						accounts[i].Distance = &distanceStr
					}
				}
			}
		}
	}

	return accounts, total, nil
}

const EarthRadius = 6371000 // meters

// Function to calculate the distance between two points on the Earth using the Haversine formula
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

// Convert degrees to radians
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func (r *accountRepository) CreateAccount(requestBody map[string]string, userID int) ([]models.Account, error) {
	account := models.Account{
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
		IsSkulid:                func(u uint) *uint { return &u }(0),
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

// Fungsi untuk menghitung account dengan filter
func (r *accountRepository) CountAccount(
	userRole string,
	territoryID int,
	withGeoJson bool,
) (int64, map[string]int64, []map[string]interface{}, response.TerritoryInfo, error) {
	var total int64
	categories := make(map[string]int64)
	var currentTerritory response.TerritoryInfo

	type CategoryCount struct {
		AccountCategory string
		Count           int64
	}

	baseQuery := r.db.Model(&models.Account{}).
		Joins("LEFT JOIN cities ON accounts.city = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id")

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

	// Total akun
	if err := baseQuery.Count(&total).Error; err != nil {
		return 0, nil, nil, currentTerritory, err
	}

	// Kategori akun
	var categoryResults []CategoryCount
	if err := baseQuery.
		Select("account_category, COUNT(*) AS count").
		Group("account_category").
		Scan(&categoryResults).Error; err != nil {
		return 0, nil, nil, currentTerritory, err
	}

	for _, res := range categoryResults {
		categories[res.AccountCategory] = res.Count
	}

	// Detail berdasarkan wilayah
	countByTerritories, err := r.CountByTerritories(userRole, territoryID, withGeoJson)
	if err != nil {
		return 0, nil, nil, currentTerritory, err
	}

	// Set wilayah saat ini
	currentTerritory, _ = r.GetTerritoryInfo(userRole, territoryID)

	return total, categories, countByTerritories, currentTerritory, nil
}

// Mengambil informasi wilayah saat ini
func (r *accountRepository) GetTerritoryInfo(role string, id int) (response.TerritoryInfo, error) {
	var info response.TerritoryInfo
	var name string
	var geojson *string

	table, err := getGeojsonTable(role)
	if err != nil {
		return info, nil // return default kosong jika tidak ada
	}

	if role == "Super-Admin" || role == "HQ" {
		return response.TerritoryInfo{ID: 0, Name: "Indonesia", Geojson: ""}, nil
	}

	err = r.db.Table(table).Select("name, geojson").Where("id = ?", id).Row().Scan(&name, &geojson)
	if err != nil {
		return info, err
	}

	gjson := ""
	if geojson != nil {
		gjson = *geojson
	}

	return response.TerritoryInfo{ID: id, Name: name, Geojson: gjson}, nil
}

func (r *accountRepository) CountByTerritories(
	userRole string,
	territoryID int,
	withGeoJson bool,
) ([]map[string]interface{}, error) {

	var results []Result
	query := r.db.Model(&models.Account{}).
		Joins("LEFT JOIN cities ON accounts.city = cities.id").
		Joins("LEFT JOIN clusters ON cities.cluster_id = clusters.id").
		Joins("LEFT JOIN branches ON clusters.branch_id = branches.id").
		Joins("LEFT JOIN regions ON branches.region_id = regions.id").
		Joins("LEFT JOIN areas ON regions.area_id = areas.id")

	var groupCol, nameCol, idCol, nextTerritory string

	switch userRole {
	case "Super-Admin", "HQ":
		groupCol, nameCol, idCol = "areas.id", "areas.name", "areas.id"
		nextTerritory = "Area"
	case "Area":
		query = query.Where("areas.id = ?", territoryID)
		groupCol, nameCol, idCol = "regions.id", "regions.name", "regions.id"
		nextTerritory = "Regional"
	case "Regional":
		query = query.Where("regions.id = ?", territoryID)
		groupCol, nameCol, idCol = "branches.id", "branches.name", "branches.id"
		nextTerritory = "Branch"
	case "Branch", "Buddies", "DS", "YAE", "Organic":
		query = query.Where("branches.id = ?", territoryID)
		groupCol, nameCol, idCol = "clusters.id", "clusters.name", "clusters.id"
		nextTerritory = "Cluster"
	case "Cluster", "Admin-Tap":
		query = query.Where("clusters.id = ?", territoryID)
		groupCol, nameCol, idCol = "cities.id", "cities.name", "cities.id"
		nextTerritory = ""
	default:
		return nil, fmt.Errorf("role %s not supported", userRole)
	}

	err := query.Select(fmt.Sprintf("%s AS territory_id, %s AS territory_name, account_category, COUNT(*) AS count",
		idCol, nameCol)).
		Group(fmt.Sprintf("%s, %s, account_category", groupCol, nameCol)).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Ambil GeoJSON sekali untuk semua wilayah
	geojsonMap := make(map[int]string)
	if withGeoJson && len(results) > 0 {
		ids := getTerritoryIDs(results)
		geojsonMap, err = r.fetchGeojsonByRole(userRole, ids)
		if err != nil {
			return nil, err
		}
	}

	// Gabungkan data
	grouped := make(map[int]map[string]interface{})
	for _, res := range results {
		if _, ok := grouped[res.TerritoryID]; !ok {
			geo := ""
			if withGeoJson {
				geo = geojsonMap[res.TerritoryID]
			}
			grouped[res.TerritoryID] = map[string]interface{}{
				"id":             res.TerritoryID,
				"name":           res.TerritoryName,
				"geojson":        geo,
				"next_territory": nextTerritory,
				"total":          int64(0),
				"categories":     map[string]int64{},
			}
		}
		entry := grouped[res.TerritoryID]
		entry["total"] = entry["total"].(int64) + res.Count
		entry["categories"].(map[string]int64)[res.AccountCategory] = res.Count
	}

	// Konversi ke slice
	final := make([]map[string]interface{}, 0, len(grouped))
	for _, val := range grouped {
		final = append(final, val)
	}
	return final, nil
}

type Result struct {
	TerritoryID     int
	TerritoryName   string
	Geojson         *string
	AccountCategory string
	Count           int64
}

func getTerritoryIDs(results []Result) []int {
	idSet := map[int]struct{}{}
	for _, r := range results {
		idSet[r.TerritoryID] = struct{}{}
	}
	ids := make([]int, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	return ids
}

func (r *accountRepository) fetchGeojsonByRole(role string, ids []int) (map[int]string, error) {
	table, err := getGeojsonTable(role)
	if err != nil {
		return nil, err
	}

	type Row struct {
		ID      int
		Geojson *string
	}

	var rows []Row
	err = r.db.Table(table).Select("id, geojson").Where("id IN ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)
	for _, row := range rows {
		if row.Geojson != nil {
			result[row.ID] = *row.Geojson
		} else {
			result[row.ID] = ""
		}
	}
	return result, nil
}

func getGeojsonTable(role string) (string, error) {
	switch role {
	case "Super-Admin", "HQ":
		return "areas", nil
	case "Area":
		return "regions", nil
	case "Regional":
		return "branches", nil
	case "Branch", "Buddies", "DS", "YAE", "Organic":
		return "clusters", nil
	case "Cluster", "Admin-Tap":
		return "cities", nil
	default:
		return "", fmt.Errorf("unsupported role: %s", role)
	}
}

func (r *accountRepository) FindAccountsWithDifferentPic(accountIDs []int, userID int) ([]models.Account, error) {
	var accounts []models.Account
	if len(accountIDs) == 0 {
		return accounts, nil
	}

	userIDStr := fmt.Sprintf("%d", userID)
	err := r.db.Model(&models.Account{}).
		Where("id IN ?", accountIDs).
		Where("pic IS NOT NULL").
		Where("pic != ''").
		Where("pic != ?", userIDStr).
		Find(&accounts).Error

	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) UpdatePicMultipleAccounts(accountIDs []int, picID int) error {
	if len(accountIDs) == 0 {
		return nil // Tidak ada akun untuk diupdate
	}

	// Update hanya kolom 'pic'
	return r.db.Model(&models.Account{}).
		Where("id IN ?", accountIDs).
		Updates(map[string]interface{}{"pic": picID}).Error
}

func (r *accountRepository) UpdateAccountsPriority(accountIDs []int, priority string) error {
	if len(accountIDs) == 0 {
		return nil
	}

	if err := r.db.Model(&models.Account{}).
		Where("id IN ?", accountIDs).
		Update("priority", priority).Error; err != nil {
		return err
	}

	return nil
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

func (r *accountRepository) FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.AccountResponse, error) {
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
		Preload("Products").
		Preload("PicDetail").
		Preload("PicInternalDetail").
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
			if account.Pic == nil || *account.Pic == "" || *account.Pic == fmt.Sprintf("%d", userID) {
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

func (r *accountRepository) GetAccountVisitCounts(
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

func (r *accountRepository) CheckAlreadyUpdateData(accountID int, userID int, clockInTime time.Time) (bool, error) {
	var count int64
	now := time.Now()

	err := r.db.
		Table("history_activity_accounts").
		Where("account_id = ? AND user_id = ? AND type = ? AND DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') >= ? AND DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') <= ?",
			accountID, userID, "Update Account", clockInTime.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05")).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *accountRepository) CreateHistoryActivityAccount(userID, accountID uint, updateType string, subjectType *string, subjectID *uint) error {
	history := models.HistoryActivityAccount{
		UserID:      userID,
		AccountID:   accountID,
		Type:        updateType,
		SubjectType: subjectType,
		SubjectID:   subjectID,
	}
	return r.db.Table("history_activity_accounts").Create(&history).Error
}
