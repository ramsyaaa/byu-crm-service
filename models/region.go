package models

import "time"

type Region struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	AreaID    *string   `json:"area_id"`
	Geojson   *string   `json:"geojson" gorm:"type:longtext"`
	Area      *Area     `json:"area" gorm:"foreignKey:AreaID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
