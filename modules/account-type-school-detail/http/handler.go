package http

import (
	"byu-crm-service/modules/account-type-school-detail/service"
)

type AccountTypeSchoolDetailHandler struct {
	service service.AccountTypeSchoolDetailService
}

func NewAccountTypeSchoolDetailHandler(service service.AccountTypeSchoolDetailService) *AccountTypeSchoolDetailHandler {
	return &AccountTypeSchoolDetailHandler{service: service}
}
