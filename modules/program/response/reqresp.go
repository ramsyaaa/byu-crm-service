package response

import "time"

type ProgramResponse struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProgramName      *string    `json:"program_name"`
	ProgramCategory  *string    `json:"program_category"`
	ProgramType      *string    `json:"program_type"`
	StartDateProject *time.Time `json:"start_date_project"`
	EndDateProject   *time.Time `json:"end_date_project"`
	KeyVisual        *string    `json:"key_visual"`
	AdditionalFile   *string    `json:"additional_file"`
	Description      *string    `json:"description"`
	QuotaValue       *float32   `json:"quota_value"`
	ValidityValue    *float32   `json:"validity_value"`
	ValidityUnit     *string    `json:"validity_unit"`
}

type SingleProductResponse struct {
	ID              uint      `json:"id"`
	ProductName     string    `json:"product_name"`
	Description     string    `json:"description"`
	ProductCategory string    `json:"product_category"`
	ProductType     string    `json:"product_type"`
	Bid             *string   `json:"bid"`
	Price           *string   `json:"price"`
	KeyVisual       *string   `json:"key_visual"`
	AdditionalFile  *string   `json:"additional_file"`
	QuotaValue      *string   `json:"quota_value"`
	ValidityValue   *string   `json:"validity_value"`
	ValidityUnit    *string   `json:"validity_unit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	EligibilityCategory []string            `json:"eligibility_category"`
	EligibilityType     []string            `json:"eligibility_type"`
	EligibilityLocation map[string][]string `json:"eligibility_location"`
}
