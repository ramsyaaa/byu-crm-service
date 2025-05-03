package repository

import (
	"byu-crm-service/models"
)

type CommunicationRepository interface {
	GetAllCommunications(limit int, paginate bool, page int, filters map[string]string, accountID int) ([]models.Communication, int64, error)
	FindByCommunicationID(id uint) (*models.Communication, error)
	CreateCommunication(requestBody map[string]string) (*models.Communication, error)
	UpdateFields(id uint, fields map[string]interface{}) error
}
