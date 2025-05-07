package models

import "time"

type AccountProduct struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID *uint     `json:"account_id"`
	ProductID *uint     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
