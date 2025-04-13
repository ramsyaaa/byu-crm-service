package models

type Region struct {
	ID     uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string  `json:"name"`
	AreaID *string `json:"area_id"`
	Area   *Area   `json:"area" gorm:"foreignKey:AreaID"`
}
