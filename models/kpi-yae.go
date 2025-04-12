package models

import "time"

type KpiYae struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
