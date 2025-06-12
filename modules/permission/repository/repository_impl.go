package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/permission/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetAllPermissions(limit int, paginate bool, page int, filters map[string]string) ([]response.PermissionResponse, int64, error) {
	var permissions []response.PermissionResponse
	var total int64

	query := r.db.Model(&models.Permission{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("permissions.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("permissions.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("permissions.created_at <= ?", endDate)
	}

	// Get total count before applying pagination
	query.Count(&total)

	if orderBy, exists := filters["order_by"]; exists && orderBy != "" {
		if order, exists := filters["order"]; exists && order != "" {
			query = query.Order(orderBy + " " + order)
		}
	}

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) GetAllPermissionsByRoleID(roleID int) ([]response.PermissionResponse, error) {
	var permissions []response.PermissionResponse

	err := r.db.
		Table("permissions").
		Select("permissions.id, permissions.name, permissions.guard_name").
		Joins("JOIN role_has_permissions ON permissions.id = role_has_permissions.permission_id").
		Where("role_has_permissions.role_id = ?", roleID).
		Scan(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepository) GetPermissionByID(id int) (*response.PermissionResponse, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}

	PermissionResponse := &response.PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		GuardName: permission.GuardName,
	}

	return PermissionResponse, nil
}

func (r *permissionRepository) GetPermissionByName(name string) (*response.PermissionResponse, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}

	permissionResponse := &response.PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		GuardName: permission.GuardName,
	}

	return permissionResponse, nil
}

func (r *permissionRepository) CreatePermission(permission *models.Permission) (*response.PermissionResponse, error) {
	if err := r.db.Create(permission).Error; err != nil {
		return nil, err
	}

	var createdPermission models.Permission
	if err := r.db.First(&createdPermission, "id = ?", permission.ID).Error; err != nil {
		return nil, err
	}

	err := r.AddRolePermissions(2, []int{int(createdPermission.ID)})
	if err != nil {
		return nil, err
	}

	permissionResponse := &response.PermissionResponse{
		ID:        createdPermission.ID,
		Name:      createdPermission.Name,
		GuardName: createdPermission.GuardName,
	}

	return permissionResponse, nil
}

func (r *permissionRepository) UpdatePermission(permission *models.Permission, id int) (*response.PermissionResponse, error) {
	var existing models.Permission

	if err := r.db.Model(&existing).Where("id = ?", id).Updates(permission).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &response.PermissionResponse{
		ID:        existing.ID,
		Name:      existing.Name,
		GuardName: existing.GuardName,
	}, nil
}

func (r *permissionRepository) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	// Begin transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Step 1: Delete existing permissions for the role
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleHasPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 2: Insert new permissions
	var rolePermissions []models.RoleHasPermission
	for _, pid := range permissionIDs {
		rolePermissions = append(rolePermissions, models.RoleHasPermission{
			RoleID:       strconv.Itoa(roleID),
			PermissionID: uint(pid),
		})
	}

	if len(rolePermissions) > 0 {
		if err := tx.Create(&rolePermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

func (r *permissionRepository) AddRolePermissions(roleID int, permissionIDs []int) error {
	// Begin transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, pid := range permissionIDs {
		// Check if the pair already exists
		var count int64
		err := tx.Model(&models.RoleHasPermission{}).
			Where("role_id = ? AND permission_id = ?", roleID, pid).
			Count(&count).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// If not exists, insert it
		if count == 0 {
			rp := models.RoleHasPermission{
				RoleID:       strconv.Itoa(roleID),
				PermissionID: uint(pid),
			}
			if err := tx.Create(&rp).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit transaction
	return tx.Commit().Error
}
