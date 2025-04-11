package http

import (
	"byu-crm-service/modules/kpi-yae-range/service"
)

type KpiYaeRangeHandler struct {
	service service.KpiYaeRangeService
}

func NewKpiYaeRangeHandler(service service.KpiYaeRangeService) *KpiYaeRangeHandler {
	return &KpiYaeRangeHandler{service: service}
}
