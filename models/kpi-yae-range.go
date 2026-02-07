package models

import (
	"time"
)

type KpiYaeRange struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Target    string     `json:"target"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserKpiPerformance struct {
	UserID       uint              `json:"user_id"`
	Name         string            `json:"name"`
	YaeCode      *string           `json:"yae_code"`
	Performances []UserPerformance `json:"performances"`
}
