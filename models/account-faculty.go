package models

import "time"

type AccountFaculty struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID *uint     `json:"account_id"`
	FacultyID *uint     `json:"faculty_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Faculty Faculty `gorm:"foreignKey:FacultyID" json:"faculty"`
}
