package models

import "time"

type ConstantData struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Type      *string   `json:"type"`
	Value     *string   `json:"value"`
	Label     *string   `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
