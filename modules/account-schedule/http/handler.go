package http

import (
	"byu-crm-service/modules/account-schedule/service"
)

type AccountScheduleHandler struct {
	service service.AccountScheduleService
}

func NewAccountScheduleHandler(service service.AccountScheduleService) *AccountScheduleHandler {
	return &AccountScheduleHandler{service: service}
}
