package models

import "time"

type Product struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Bid             *string   `json:"bid"`
	Price           *string   `json:"price"`
	ProductName     *string   `json:"product_name"`
	ProductCategory *string   `json:"product_category"`
	ProductType     *string   `json:"product_type"`
	KeyVisual       *string   `json:"key_visual"`
	AdditionalFile  *string   `json:"additional_file"`
	Description     *string   `json:"description"`
	QuotaValue      *float32  `json:"quota_value"`
	ValidityValue   *float32  `json:"validity_value"`
	ValidityUnit    *string   `json:"validity_unit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Eligibility *Eligibility `json:"eligibility" gorm:"foreignKey:SubjectID;references:ID;where:subject_type='App\\Models\\Product'"`
}
