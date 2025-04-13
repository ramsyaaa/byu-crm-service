package models

type Cluster struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string  `json:"name"`
	BranchID *string `json:"branch_id"`
	Branch   *Branch `json:"branch" gorm:"foreignKey:BranchID"`
}
