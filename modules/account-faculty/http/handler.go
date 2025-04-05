package http

import (
	"byu-crm-service/modules/account-faculty/service"
)

type AccountFacultyHandler struct {
	service service.AccountFacultyService
}

func NewAccountFacultyHandler(service service.AccountFacultyService) *AccountFacultyHandler {
	return &AccountFacultyHandler{service: service}
}
