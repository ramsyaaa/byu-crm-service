package models

import (
	"time"
)

type Contact struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ContactName *string    `json:"contact_name"`
	PhoneNumber *string    `json:"phone_number"`
	Position    *string    `json:"position"`
	Birthday    *time.Time `json:"birthday"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Accounts []Account `gorm:"many2many:contact_accounts;foreignKey:ID;joinForeignKey:contact_id;References:ID;joinReferences:account_id" json:"accounts"`
}
