package models

import "time"

type DetailCommunityMember struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID    uint       `json:"account_id"`
	Name         *string    `json:"name"`
	Gender       *string    `json:"gender"`
	City         *string    `json:"city"`
	Phone        *string    `json:"phone"`
	JoinDate     *time.Time `json:"join_date"`
	UploadedDate *time.Time `json:"uploaded_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
