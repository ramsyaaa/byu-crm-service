package models

import (
	"time"
)

// ApiLog represents the api_logs table for storing API request/response logs
type ApiLog struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccessedAt       time.Time `gorm:"not null;index" json:"accessed_at"`
	Endpoint         string    `gorm:"type:text;not null" json:"endpoint"`
	Method           string    `gorm:"type:varchar(10);not null;index" json:"method"`
	StatusCode       int       `gorm:"not null;index" json:"status_code"`
	ResponseTimeMs   int64     `gorm:"not null" json:"response_time_ms"`
	RequestPayload   *string   `gorm:"type:longtext" json:"request_payload"`
	ResponsePayload  *string   `gorm:"type:longtext" json:"response_payload"`
	ErrorMessage     *string   `gorm:"type:text" json:"error_message"`
	AuthUserEmail    *string   `gorm:"type:varchar(255);index" json:"auth_user_email"`
	IPAddress        string    `gorm:"type:varchar(45);not null;index" json:"ip_address"`
	UserAgent        string    `gorm:"type:text" json:"user_agent"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for the ApiLog model
func (ApiLog) TableName() string {
	return "api_logs"
}
