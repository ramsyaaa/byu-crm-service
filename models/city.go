package models

import "time"

type City struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	ClusterID *string   `json:"cluster_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cluster   *Cluster  `json:"cluster" gorm:"foreignKey:ClusterID"`
}
