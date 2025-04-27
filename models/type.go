package models

import "time"

type Type struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  *uint     `json:"category_id"`
	ModuleType  *string   `json:"module_type"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
