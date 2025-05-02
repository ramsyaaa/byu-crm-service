package models

import (
	"time"
)

type Communication struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CommunicationType *string    `json:"communication_type"`
	AccountID         *uint      `json:"account_id"`
	ContactID         *uint      `json:"contact_id"`
	Note              *string    `json:"note"`
	Date              *time.Time `json:"date"`
	CreatedBy         *uint      `json:"created_by"`

	MainCommunicationID *uint          `json:"main_communication_id"`
	MainCommunication   *Communication `gorm:"foreignKey:MainCommunicationID" json:"main_communication"`

	NextCommunicationID *uint          `json:"next_communication_id"`
	NextCommunication   *Communication `gorm:"foreignKey:NextCommunicationID" json:"next_communication"`

	Account *Account `gorm:"foreignKey:AccountID" json:"account"`
	Contact *Contact `gorm:"foreignKey:ContactID" json:"contact"`

	StatusCommunication *string `json:"status_communication"`

	OpportunityID *uint        `json:"opportunity_id"`
	Opportunity   *Opportunity `gorm:"foreignKey:OpportunityID" json:"opportunity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
