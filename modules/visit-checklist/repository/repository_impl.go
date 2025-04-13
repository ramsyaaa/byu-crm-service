package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type visitChecklistRepository struct {
	db *gorm.DB
}

func NewVisitChecklistRepository(db *gorm.DB) VisitChecklistRepository {
	return &visitChecklistRepository{db: db}
}

func (r *visitChecklistRepository) GetAllVisitChecklist() ([]models.VisitChecklist, error) {
	var checklists []models.VisitChecklist

	err := r.db.Find(&checklists).Error
	return checklists, err
}
