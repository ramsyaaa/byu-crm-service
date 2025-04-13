package http

import (
	"byu-crm-service/modules/visit-history/service"
)

type VisitHistoryHandler struct {
	visitHistoryService service.VisitHistoryService
}

func NewVisitHistoryHandler(visitHistoryService service.VisitHistoryService) *VisitHistoryHandler {
	return &VisitHistoryHandler{visitHistoryService: visitHistoryService}
}
