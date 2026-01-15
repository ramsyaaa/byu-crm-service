package repository

import (
	"byu-crm-service/models"
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type performanceDigiposRepository struct {
	db *gorm.DB
}

func NewPerformanceDigiposRepository(db *gorm.DB) PerformanceDigiposRepository {
	return &performanceDigiposRepository{db: db}
}

func (r *performanceDigiposRepository) Create(performanceDigipos *models.PerformanceDigipos) error {
	return r.db.Create(performanceDigipos).Error
}

func (r *performanceDigiposRepository) FindByIdImport(idImport string) (*models.PerformanceDigipos, error) {
	var performance models.PerformanceDigipos
	err := r.db.Where("id_import = ?", idImport).First(&performance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Data tidak ditemukan, kembalikan nil
		}
		return nil, err
	}
	return &performance, nil
}

func (r *performanceDigiposRepository) Update(performance *models.PerformanceDigipos) error {
	return r.db.Save(performance).Error
}

func (r *performanceDigiposRepository) CountPerformanceByUserYaeCode(
	userID int,
	month uint,
	year uint,
) (int, error) {

	var count int64
	var yaeCode sql.NullString

	// 1. Ambil yae_code dari users
	err := r.db.
		Table("users").
		Select("yae_code").
		Where("id = ?", userID).
		Scan(&yaeCode).Error
	if err != nil {
		return 0, err
	}

	// 2. Jika NULL atau string kosong → return 0
	if !yaeCode.Valid || strings.TrimSpace(yaeCode.String) == "" {
		return 0, nil
	}

	// 3. Query ke performances
	query := r.db.
		Model(&models.PerformanceDigipos{}).
		Where("event_name LIKE ?", "%"+yaeCode.String+"%")

	// 4. Filter bulan & tahun (opsional)
	if month != 0 && year != 0 {
		query = query.Where(
			"MONTH(created_at) = ? AND YEAR(created_at) = ?",
			month,
			year,
		)
	}

	// 5. Hitung
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
