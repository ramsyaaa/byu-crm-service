package models

import "time"

type Eligibility struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectType string    `json:"subject_type"`
	SubjectID   uint      `json:"subject_id"`
	Categories  string    `json:"categories"`
	Types       string    `json:"types"`
	Locations   string    `json:"locations"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
