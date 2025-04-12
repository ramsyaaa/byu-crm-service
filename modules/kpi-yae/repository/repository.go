package repository

import "byu-crm-service/models"

type KpiYaeRepository interface {
	GetKpiYaeByName(name string) (*models.KpiYae, error)
}
