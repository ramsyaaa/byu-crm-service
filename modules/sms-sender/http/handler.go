package http

import (
	"byu-crm-service/modules/sms-sender/service"
)

type SmsSenderHandler struct {
	smsSenderService service.SmsSenderService
}

func NewSmsSenderHandler(smsSenderService service.SmsSenderService) *SmsSenderHandler {
	return &SmsSenderHandler{smsSenderService: smsSenderService}
}
