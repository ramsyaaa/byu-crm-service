package models

import "time"

type UserStatusHistory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    *uint     `gorm:"column:user_id" json:"user_id"`
	StartDate time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   time.Time `gorm:"column:end_date" json:"end_date"`
	Status    *string   `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
