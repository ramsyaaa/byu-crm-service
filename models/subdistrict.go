package models

import "time"

type Subdistrict struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	CityID    *string   `json:"city_id"`
	City      *City     `json:"city" gorm:"foreignKey:CityID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
