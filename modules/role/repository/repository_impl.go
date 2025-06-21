package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/role/response"
	"strings"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAllRoles(limit int, paginate bool, page int, filters map[string]string) ([]response.RoleResponse, int64, error) {
	var roles []response.RoleResponse
	var total int64

	query := r.db.Model(&models.Role{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("roles.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("roles.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("roles.created_at <= ?", endDate)
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

	err := query.Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) AssignModelHasRole(model_type string, model_id int, role_id int) error {
	if err := r.db.
		Table("model_has_roles").
		Where("model_type = ? AND model_id = ?", model_type, model_id).
		Delete(nil).Error; err != nil {
		return err
	}

	// Buat map untuk data yang akan dimasukkan
	data := map[string]interface{}{
		"role_id":    role_id,
		"model_type": model_type,
		"model_id":   model_id,
	}

	// Masukkan ke tabel model_has_roles
	if err := r.db.Table("model_has_roles").Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *roleRepository) GetRoleByID(id int) (*response.RoleResponse, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}

	RoleResponse := &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		GuardName: role.GuardName,
	}

	return RoleResponse, nil
}

func (r *roleRepository) GetRoleByName(name string) (*response.RoleResponse, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}

	RoleResponse := &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		GuardName: role.GuardName,
	}

	return RoleResponse, nil
}

func (r *roleRepository) CreateRole(role *models.Role) (*response.RoleResponse, error) {
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}

	var createdRole models.Role
	if err := r.db.First(&createdRole, "id = ?", role.ID).Error; err != nil {
		return nil, err
	}

	RoleResponse := &response.RoleResponse{
		ID:        createdRole.ID,
		Name:      createdRole.Name,
		GuardName: createdRole.GuardName,
	}

	return RoleResponse, nil
}

func (r *roleRepository) UpdateRole(role *models.Role, id int) (*response.RoleResponse, error) {
	var existing models.Role

	if err := r.db.Model(&existing).Where("id = ?", id).Updates(role).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &response.RoleResponse{
		ID:        existing.ID,
		Name:      existing.Name,
		GuardName: existing.GuardName,
	}, nil
}
