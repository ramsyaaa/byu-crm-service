package response

import (
	"byu-crm-service/models"
	"time"
)

type ContactResponse struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ContactName *string    `json:"contact_name"`
	PhoneNumber *string    `json:"phone_number"`
	Position    *string    `json:"position"`
	Birthday    *time.Time `json:"birthday"`

	Accounts []models.Account `json:"accounts"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
