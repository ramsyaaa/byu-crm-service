package http

import (
	"byu-crm-service/modules/account-type-campus-detail/service"
)

type AccountTypeCampusDetailHandler struct {
	service service.AccountTypeCampusDetailService
}

func NewAccountTypeCampusDetailHandler(service service.AccountTypeCampusDetailService) *AccountTypeCampusDetailHandler {
	return &AccountTypeCampusDetailHandler{service: service}
}
