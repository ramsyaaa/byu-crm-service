package models

type City struct {
	ID        uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string   `json:"name"`
	ClusterID *string  `json:"cluster_id"`
	Cluster   *Cluster `json:"cluster" gorm:"foreignKey:ClusterID"`
}
