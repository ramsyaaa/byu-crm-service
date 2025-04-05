package http

import (
	"byu-crm-service/modules/contact-account/service"
)

type ContactAccountHandler struct {
	service service.ContactAccountService
}

func NewContactAccountHandler(service service.ContactAccountService) *ContactAccountHandler {
	return &ContactAccountHandler{service: service}
}
