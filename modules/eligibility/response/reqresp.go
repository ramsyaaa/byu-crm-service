package response

type Eligibility struct {
	ID          uint   `gorm:"primaryKey"`
	SubjectType string `gorm:"column:subject_type"`
	SubjectID   int    `gorm:"column:subject_id"`
	Categories  string `gorm:"column:categories"` // JSON string
	Types       string `gorm:"column:types"`      // JSON string
	Locations   string `gorm:"column:locations"`  // JSON string
}

type LocationFilter struct {
	Areas     []string `json:"areas"`
	Regions   []string `json:"regions"`
	Branches  []string `json:"branches"`
	Clusters  []string `json:"clusters"`
	Cities    []string `json:"cities"`
	Districts []string `json:"districts"`
}
