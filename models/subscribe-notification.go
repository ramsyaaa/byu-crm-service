package models

import "time"

type SubscribeNotification struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      *uint     `json:"user_id"`
	SubscribeID string    `json:"subscribe_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
