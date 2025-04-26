package models

import "time"

type Branch struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	RegionID  *string   `json:"region_id"`
	Region    *Region   `json:"region" gorm:"foreignKey:RegionID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
