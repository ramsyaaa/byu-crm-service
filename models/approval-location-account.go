package models

import "time"

type ApprovalLocationAccount struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint     `json:"user_id"`
	AccountID *uint     `json:"account_id"`
	Longitude *string   `json:"longitude"`
	Latitude  *string   `json:"latitude"`
	Status    *uint     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
