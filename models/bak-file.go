package models

import (
	"time"
)

type BakFile struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProgramName  string    `json:"program_name"`
	ContractDate time.Time `json:"contract_date"`
	AccountID    uint      `json:"account_id"`
	UserID       uint      `json:"user_id"`

	// First party
	FirstPartyName        string `json:"first_party_name"`
	FirstPartyPosition    string `json:"first_party_position"`
	FirstPartyPhoneNumber string `json:"first_party_phone_number"`
	FirstPartyAddress     string `json:"first_party_address"`

	// Second party
	SecondPartyCompany     string `json:"second_party_company"`
	SecondPartyName        string `json:"second_party_name"`
	SecondPartyPosition    string `json:"second_party_position"`
	SecondPartyPhoneNumber string `json:"second_party_phone_number"`
	SecondPartyAddress     string `json:"second_party_address"`

	// Description
	Description string `json:"description"`

	// Additional signs (multiple values, comma separated or stored in JSONB in DB if needed)
	AdditionalSignTitle    *string `json:"additional_sign_title,omitempty"` // use *string if nullable
	AdditionalSignName     *string `json:"additional_sign_name,omitempty"`
	AdditionalSignPosition *string `json:"additional_sign_position,omitempty"`

	Account Account `gorm:"foreignKey:AccountID" json:"account"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
