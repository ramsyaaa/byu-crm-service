package models

import "time"

type AccountSchedule struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectType      *string    `json:"subject_type"`
	SubjectID        *uint      `json:"subject_id"`
	Date             *time.Time `json:"date"`
	Title            *string    `json:"title"`
	ScheduleCategory *string    `json:"schedule_category"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
