package models

type PerformanceIndiana struct {
	ID             uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         *int    `gorm:"column:user_id" json:"user_id"`
	DataIn         *int    `gorm:"column:data_in" json:"data_in"`
	NotProcess     *int    `gorm:"column:not_process" json:"not_process"`
	Rejected       *int    `gorm:"column:rejected" json:"rejected"`
	Pending        *int    `gorm:"column:pending" json:"pending"`
	Approve        *int    `gorm:"column:approve" json:"approve"`
	InProgress     *int    `gorm:"column:in_progress" json:"in_progress"`
	Failed         *int    `gorm:"column:failed" json:"failed"`
	Active         *int    `gorm:"column:active" json:"active"`
	ActiveBySystem *int    `gorm:"column:active_by_system" json:"active_by_system"`
	Month          *string `gorm:"column:month" json:"month"`
}

func (PerformanceIndiana) TableName() string {
	return "performance_indianas"
}
