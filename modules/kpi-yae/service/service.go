package service

import "byu-crm-service/models"

type KpiYaeService interface {
	GetKpiYaeByName(name string) (*models.KpiYae, error)
}
