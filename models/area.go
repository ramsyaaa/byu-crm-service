package models

type Area struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}
