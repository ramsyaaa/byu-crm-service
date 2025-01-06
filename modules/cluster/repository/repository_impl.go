package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type clusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db}
}

func (r *clusterRepository) FindByID(id uint) (*models.Cluster, error) {
	var cluster models.Cluster
	if err := r.db.First(&cluster, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &cluster, nil
}

func (r *clusterRepository) FindByName(name string) (*models.Cluster, error) {
	var cluster models.Cluster
	if err := r.db.Where("name = ?", name).First(&cluster).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &cluster, nil
}
