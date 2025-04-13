package repository

import (
	"byu-crm-service/models"
	"fmt"

	"gorm.io/gorm"
)

type visitHistoryRepository struct {
	db *gorm.DB
}

func NewVisitHistoryRepository(db *gorm.DB) VisitHistoryRepository {
	return &visitHistoryRepository{db: db}
}

func (r *visitHistoryRepository) CreateVisitHistory(visit_history *models.VisitHistory) (*models.VisitHistory, error) {
	if err := r.db.Create(visit_history).Error; err != nil {
		return nil, err
	}

	var createdVisitHistory models.VisitHistory
	if err := r.db.First(&createdVisitHistory, "id = ?", visit_history.ID).Error; err != nil {
		return nil, err
	}

	return &createdVisitHistory, nil
}

func (r *visitHistoryRepository) CountVisitHistory(user_id int, month uint, year uint, kpi_name string) (int, error) {
	var count int64

	query := r.db.
		Model(&models.VisitHistory{}).
		Select("COUNT(DISTINCT subject_id)").
		Where("user_id = ? AND subject_type = ? AND MONTH(created_at) = ? AND YEAR(created_at) = ?",
			user_id, "App\\Models\\Account", month, year)

	if kpi_name != "" {
		// LIKE tanpa escape tanda kutip dua
		likePattern := fmt.Sprintf("%%%s\":1%%", kpi_name)
		query = query.Where("target LIKE ?", likePattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
