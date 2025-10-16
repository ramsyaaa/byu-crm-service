package models

import "time"

type ParameterCalculatorKpi struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Parameters string     `json:"parameters"`
	Incentive  string     `json:"incentive"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
