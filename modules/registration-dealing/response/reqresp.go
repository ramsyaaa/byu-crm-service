package response

import (
	"time"
)

type RegistrationDealingResponse struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber          *string   `json:"phone_number"`
	AccountID            *uint     `json:"account_id"`
	CustomerName         *string   `json:"customer_name"`
	EventName            *string   `json:"event_name"`
	RegistrationEvidence *string   `json:"registration_evidence"`
	WhatsappNumber       *string   `json:"whatsapp_number"`
	Class                *string   `json:"class"`
	Email                *string   `json:"email"`
	UserID               *uint     `json:"user_id"`
	Source               *string   `json:"source"`
	SchoolType           *string   `json:"school_type"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
