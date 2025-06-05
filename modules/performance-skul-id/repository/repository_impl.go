package repository

import (
	"byu-crm-service/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type performanceSkulIdRepository struct {
	db *gorm.DB
}

func NewPerformanceSkulIdRepository(db *gorm.DB) PerformanceSkulIdRepository {
	return &performanceSkulIdRepository{db: db}
}

func (r *performanceSkulIdRepository) Create(performanceSkulId *models.PerformanceSkulId) error {
	return r.db.Create(performanceSkulId).Error
}

func (r *performanceSkulIdRepository) FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error) {
	var performanceSkulId models.PerformanceSkulId
	err := r.db.Where("msisdn = ?", serial).First(&performanceSkulId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &performanceSkulId, nil
}

func (r *performanceSkulIdRepository) FindByIdSkulId(idSkulId string) (*models.PerformanceSkulId, error) {
	var performance models.PerformanceSkulId
	err := r.db.Where("id_skulid = ?", idSkulId).First(&performance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Data tidak ditemukan, kembalikan nil
		}
		return nil, err
	}
	return &performance, nil
}

func (r *performanceSkulIdRepository) FindAll(limit, offset int, filters map[string]string, accountID int, page int, paginate bool) ([]models.PerformanceSkulId, int64, error) {
	var skulids []models.PerformanceSkulId
	var total int64

	query := r.db.Model(&models.PerformanceSkulId{}).
		Select("id_skulid, user_name, msisdn, provider, registered_date").
		Where("account_id = ?", accountID)

	// Optional: filter user_type jika tidak kosong
	if filters["user_type"] != "" {
		query = query.Where("user_type = ?", filters["user_type"])
	}

	// Hitung total awal
	query.Count(&total)

	// Search filter
	if search := filters["search"]; search != "" {
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where(r.db.
			Where("LOWER(id_skulid) LIKE ?", search).
			Or("LOWER(user_name) LIKE ?", search).
			Or("LOWER(msisdn) LIKE ?", search).
			Or("LOWER(provider) LIKE ?", search))
	}

	// Order by & direction
	orderBy := filters["order_by"]
	if orderBy == "" {
		orderBy = "id"
	}

	order := filters["order"]
	if order == "" {
		order = "DESC"
	}

	query = query.Order(orderBy + " " + order)

	// Handle paginate
	if paginate {
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	// Eksekusi query
	err := query.Find(&skulids).Error
	if err != nil {
		return nil, 0, err
	}

	return skulids, total, nil
}

func (r *performanceSkulIdRepository) Update(performance *models.PerformanceSkulId) error {
	return r.db.Save(performance).Error
}
