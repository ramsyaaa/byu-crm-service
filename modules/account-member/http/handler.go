package http

import (
	"byu-crm-service/modules/account-member/service"
)

type AccountMemberHandler struct {
	service service.AccountMemberService
}

func NewAccountMemberHandler(service service.AccountMemberService) *AccountMemberHandler {
	return &AccountMemberHandler{service: service}
}
