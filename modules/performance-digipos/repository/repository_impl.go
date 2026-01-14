package repository

import (
	"byu-crm-service/models"
	"errors"

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
	var yaeCode string

	// 1. Ambil yae_code dari table users
	err := r.db.
		Table("users").
		Select("yae_code").
		Where("id = ?", userID).
		Scan(&yaeCode).Error

	if err != nil {
		return 0, err
	}

	if yaeCode == "" {
		return 0, nil // tidak ada yae_code → tidak ada data
	}

	// 2. Query ke table performances
	query := r.db.
		Model(&models.PerformanceDigipos{}).
		Where("event_name LIKE ?", "%"+yaeCode+"%")

	// 3. Filter bulan & tahun (opsional)
	if month != 0 && year != 0 {
		query = query.Where(
			"MONTH(created_at) = ? AND YEAR(created_at) = ?",
			month,
			year,
		)
	}

	// 4. Hitung jumlah data
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
