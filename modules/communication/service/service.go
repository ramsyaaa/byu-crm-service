package service

import (
	"byu-crm-service/models"
)

type CommunicationService interface {
	GetAllCommunications(limit int, paginate bool, page int, filters map[string]string, accountID int) ([]models.Communication, int64, error)
	FindByCommunicationID(id uint) (*models.Communication, error)
	CreateCommunication(requestBody map[string]interface{}, userID int) (*models.Communication, error)
	UpdateCommunication(requestBody map[string]interface{}, userID int, communicationID int) (*models.Communication, error)
}
