package http

import (
	"byu-crm-service/modules/kpi-yae/service"
)

type KpiYaeHandler struct {
	service service.KpiYaeService
}

func NewKpiYaeHandler(service service.KpiYaeService) *KpiYaeHandler {
	return &KpiYaeHandler{service: service}
}
