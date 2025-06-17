package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/user/response"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	onlyRole []string,
	orderByMostAssignedPic bool,
	userRole string,
	territoryID interface{},
) ([]response.UserResponse, int64, error) {
	var users []models.User
	var total int64

	// ======================= Base Query for Filter =======================
	baseQuery := r.db.Model(&models.User{}).
		Joins("LEFT JOIN model_has_roles ON model_has_roles.model_id = users.id AND model_has_roles.model_type = ?", "App\\Models\\User").
		Joins("LEFT JOIN roles ON roles.id = model_has_roles.role_id").
		Group("users.id")

	// Filter by Role
	if len(onlyRole) > 0 {
		baseQuery = baseQuery.Where("roles.name IN ?", onlyRole)
	}

	// Search filter
	if search, exists := filters["search"]; exists && search != "" {
		tokens := strings.Fields(search)
		for _, token := range tokens {
			baseQuery = baseQuery.Where("users.name LIKE ?", "%"+token+"%")
		}
	}

	// Date filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		baseQuery = baseQuery.Where("users.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		baseQuery = baseQuery.Where("users.created_at <= ?", endDate)
	}

	// ======================= Role-based filtering =======================

	if userRole != "" && territoryID != nil {
		switch userRole {
		case "Area":
			var regionIDs []uint
			if ids, ok := territoryID.([]uint); ok {
				r.db.Model(&models.Region{}).Where("area_id IN ?", ids).Pluck("id", &regionIDs)
			} else {
				r.db.Model(&models.Region{}).Where("area_id = ?", territoryID).Pluck("id", &regionIDs)
			}
			var branchIDs []uint
			r.db.Model(&models.Branch{}).Where("region_id IN ?", regionIDs).Pluck("id", &branchIDs)
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id IN ?", branchIDs).Pluck("id", &clusterIDs)

			baseQuery = baseQuery.Where(
				r.db.Where("territory_type = ? AND territory_id = ?", "App\\Models\\Area", territoryID).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Region", regionIDs).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Branch", branchIDs).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Region":
			var branchIDs []uint
			if ids, ok := territoryID.([]uint); ok {
				r.db.Model(&models.Branch{}).Where("region_id IN ?", ids).Pluck("id", &branchIDs)
			} else {
				r.db.Model(&models.Branch{}).Where("region_id = ?", territoryID).Pluck("id", &branchIDs)
			}
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id IN ?", branchIDs).Pluck("id", &clusterIDs)

			baseQuery = baseQuery.Where(
				r.db.Where("territory_type = ? AND territory_id = ?", "App\\Models\\Region", territoryID).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Branch", branchIDs).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Branch", "DS", "Buddies", "YAE":
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id = ?", territoryID).Pluck("id", &clusterIDs)

			baseQuery = baseQuery.Where(
				r.db.Where("territory_type = ? AND territory_id = ?", "App\\Models\\Branch", territoryID).
					Or("territory_type = ? AND territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Cluster":
			baseQuery = baseQuery.Where("territory_type = ? AND territory_id = ?", "App\\Models\\Cluster", territoryID)
		}
	}

	// ======================= Count Query =======================
	countQuery := baseQuery.Session(&gorm.Session{}) // clone
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ======================= Pagination =======================
	if paginate {
		offset := (page - 1) * limit
		baseQuery = baseQuery.Limit(limit).Offset(offset)
	} else if limit > 0 {
		baseQuery = baseQuery.Limit(limit)
	}

	// ======================= Order =======================
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderByMostAssignedPic {
		// Diambil terpisah nanti
	} else if orderBy != "" && order != "" {
		baseQuery = baseQuery.Order(orderBy + " " + order)
	}

	// ======================= Query Data Users =======================
	if err := baseQuery.Select("users.id, users.name, users.email, users.avatar, users.msisdn, users.user_status, users.user_type, users.territory_type, users.territory_id").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// ======================= Ambil TotalPic =======================
	var userIDs []uint
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	var userPics []struct {
		UserID   uint
		TotalPic int64
	}
	if len(userIDs) > 0 {
		r.db.Model(&models.Account{}).
			Select("pic as user_id, COUNT(*) as total_pic").
			Where("pic IN ?", userIDs).
			Group("pic").
			Scan(&userPics)
	}
	picMap := map[uint]int64{}
	for _, up := range userPics {
		picMap[up.UserID] = up.TotalPic
	}

	// ======================= Ambil Role =======================
	var userRoles []struct {
		UserID uint
		Name   string
	}
	if len(userIDs) > 0 {
		r.db.
			Table("roles").
			Select("model_has_roles.model_id as user_id, roles.name").
			Joins("JOIN model_has_roles ON model_has_roles.role_id = roles.id").
			Where("model_has_roles.model_id IN ? AND model_has_roles.model_type = ?", userIDs, "App\\Models\\User").
			Scan(&userRoles)
	}
	roleMap := make(map[uint][]string)
	for _, ur := range userRoles {
		roleMap[ur.UserID] = append(roleMap[ur.UserID], ur.Name)
	}

	// ======================= Bangun Response =======================
	var responses []response.UserResponse
	for _, user := range users {
		var totalPic *uint
		if v, ok := picMap[user.ID]; ok {
			u := uint(v)
			totalPic = &u
		}
		responses = append(responses, response.UserResponse{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Avatar:        user.Avatar,
			Msisdn:        user.Msisdn,
			UserStatus:    user.UserStatus,
			UserType:      user.UserType,
			TerritoryID:   user.TerritoryID,
			TerritoryType: user.TerritoryType,
			TotalPic:      totalPic,
			RoleNames:     roleMap[user.ID],
		})
	}

	return responses, total, nil
}

func (r *userRepository) FindByID(id uint) (*response.UserResponse, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Ambil role IDs dari model_has_roles
	var roleIDs []uint
	if err := r.db.Table("model_has_roles").
		Where("model_id = ? AND model_type = ?", id, "App\\Models\\User").
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}

	// Ambil nama role
	var roleNames []string
	if len(roleIDs) > 0 {
		if err := r.db.Table("roles").
			Where("id IN ?", roleIDs).
			Pluck("name", &roleNames).Error; err != nil {
			return nil, err
		}
	}

	// Ambil permission_id dari role_has_permissions
	var permissionIDs []uint
	if len(roleIDs) > 0 {
		if err := r.db.Table("role_has_permissions").
			Where("role_id IN ?", roleIDs).
			Pluck("permission_id", &permissionIDs).Error; err != nil {
			return nil, err
		}
	}

	// Ambil nama permission
	var permissions []string
	if len(permissionIDs) > 0 {
		if err := r.db.Table("permissions").
			Where("id IN ?", permissionIDs).
			Pluck("name", &permissions).Error; err != nil {
			return nil, err
		}
	}

	// Bangun response
	response := &response.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Msisdn:        user.Msisdn,
		UserStatus:    user.UserStatus,
		UserType:      user.UserType,
		TerritoryID:   user.TerritoryID,
		TerritoryType: user.TerritoryType,
		RoleNames:     roleNames,
		Permissions:   permissions,
	}

	return response, nil
}

func (r *userRepository) UpdateUserProfile(id uint, user map[string]interface{}) (*response.UserResponse, error) {
	// Ambil user yang akan diupdate
	var existingUser models.User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Siapkan map data update
	updateData := map[string]interface{}{}

	// Update nama jika tersedia
	if name, ok := user["name"].(string); ok && name != "" {
		updateData["name"] = name
	}

	// Update password jika tersedia dan tidak kosong
	if password, ok := user["password"].(string); ok && password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updateData["password"] = string(hashedPassword)
	}

	// Kalau tidak ada yang diupdate, skip
	if len(updateData) == 0 {
		return nil, nil
	}

	// Jalankan update hanya untuk field yang diperlukan
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Ambil ulang data user setelah update
	if err := r.db.First(&existingUser, id).Error; err != nil {
		return nil, err
	}

	// Buat response
	response := &response.UserResponse{
		ID:            existingUser.ID,
		Name:          existingUser.Name,
		Email:         existingUser.Email,
		Avatar:        existingUser.Avatar,
		Msisdn:        existingUser.Msisdn,
		UserStatus:    existingUser.UserStatus,
		UserType:      existingUser.UserType,
		TerritoryID:   existingUser.TerritoryID,
		TerritoryType: existingUser.TerritoryType,
	}

	return response, nil
}

func (r *userRepository) GetUserCountByRoles(
	onlyRoles []string,
	userRole string,
	territoryID interface{},
) (map[string]int64, error) {
	type RoleCount struct {
		RoleName string
		Total    int64
	}

	var results []RoleCount
	query := r.db.Table("users").
		Select("roles.name as role_name, COUNT(users.id) as total").
		Joins("JOIN model_has_roles ON model_has_roles.model_id = users.id AND model_has_roles.model_type = ?", "App\\Models\\User").
		Joins("JOIN roles ON roles.id = model_has_roles.role_id").
		Group("roles.name")

	// Filter role
	if len(onlyRoles) > 0 {
		query = query.Where("roles.name IN ?", onlyRoles)
	}

	// Filter territory
	if userRole != "" && territoryID != nil {
		switch userRole {
		case "Area":
			var regionIDs []uint
			if ids, ok := territoryID.([]uint); ok {
				r.db.Model(&models.Region{}).Where("area_id IN ?", ids).Pluck("id", &regionIDs)
			} else {
				r.db.Model(&models.Region{}).Where("area_id = ?", territoryID).Pluck("id", &regionIDs)
			}
			var branchIDs []uint
			r.db.Model(&models.Branch{}).Where("region_id IN ?", regionIDs).Pluck("id", &branchIDs)
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id IN ?", branchIDs).Pluck("id", &clusterIDs)

			query = query.Where(
				r.db.Where("users.territory_type = ? AND users.territory_id = ?", "App\\Models\\Area", territoryID).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Region", regionIDs).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Branch", branchIDs).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Region":
			var branchIDs []uint
			if ids, ok := territoryID.([]uint); ok {
				r.db.Model(&models.Branch{}).Where("region_id IN ?", ids).Pluck("id", &branchIDs)
			} else {
				r.db.Model(&models.Branch{}).Where("region_id = ?", territoryID).Pluck("id", &branchIDs)
			}
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id IN ?", branchIDs).Pluck("id", &clusterIDs)

			query = query.Where(
				r.db.Where("users.territory_type = ? AND users.territory_id = ?", "App\\Models\\Region", territoryID).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Branch", branchIDs).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Branch", "DS", "Buddies", "YAE":
			var clusterIDs []uint
			r.db.Model(&models.Cluster{}).Where("branch_id = ?", territoryID).Pluck("id", &clusterIDs)

			query = query.Where(
				r.db.Where("users.territory_type = ? AND users.territory_id = ?", "App\\Models\\Branch", territoryID).
					Or("users.territory_type = ? AND users.territory_id IN ?", "App\\Models\\Cluster", clusterIDs),
			)

		case "Cluster":
			query = query.Where("users.territory_type = ? AND users.territory_id = ?", "App\\Models\\Cluster", territoryID)
		}
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Convert to map
	countMap := make(map[string]int64)
	for _, row := range results {
		countMap[row.RoleName] = row.Total
	}

	return countMap, nil
}
