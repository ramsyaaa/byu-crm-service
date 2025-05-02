package models

import "time"

type Opportunity struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OpportunityName *string    `json:"opportunity_name"`
	AccountID       *uint      `json:"account_id"`
	ContactID       *uint      `json:"contact_id"`
	Description     *string    `json:"description"`
	OpenDate        *time.Time `json:"open_date"`
	CloseDate       *time.Time `json:"close_date"`
	Amount          *string    `json:"amount"`
	CreatedBy       *uint      `json:"created_by"`

	Account *Account `gorm:"foreignKey:AccountID" json:"account"`
	Contact *Contact `gorm:"foreignKey:ContactID" json:"contact"`
	User    *User    `gorm:"foreignKey:CreatedBy" json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
