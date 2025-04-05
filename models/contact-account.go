package models

import "time"

type ContactAccount struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID *uint     `json:"account_id"`
	ContactID *uint     `json:"contact_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
