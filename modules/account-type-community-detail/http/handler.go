package http

import (
	"byu-crm-service/modules/account-type-community-detail/service"
)

type AccountTypeCommunityDetailHandler struct {
	service service.AccountTypeCommunityDetailService
}

func NewAccountTypeCommunityDetailHandler(service service.AccountTypeCommunityDetailService) *AccountTypeCommunityDetailHandler {
	return &AccountTypeCommunityDetailHandler{service: service}
}
