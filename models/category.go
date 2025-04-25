package models

import "time"

type Category struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ModuleType *string   `json:"module_type"`
	Name       *string   `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
