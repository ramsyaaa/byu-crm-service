package models

import "time"

type Faculty struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      *string   `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
