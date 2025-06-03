package models

import "time"

type Cluster struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	BranchID  *string   `json:"branch_id"`
	Geojson   *string   `json:"geojson" gorm:"type:longtext"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Branch    *Branch   `json:"branch" gorm:"foreignKey:BranchID"`
}
