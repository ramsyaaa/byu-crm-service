package models

type User struct {
	ID                uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Avatar            string `json:"avatar"`
	Msisdn            string `json:"msisdn"`
	UserStatus        string `json:"user_status"`
	UserType          string `json:"user_type"`
	TerritoryID       uint   `json:"territory_id"`
	TerritoryType     string `json:"territory_type"`
	EmailVerifiedAt   string `json:"email_verified_at"`
	Password          string `json:"password"`
	IsRequestPassword bool   `json:"is_request_password"`
	RememberToken     string `json:"remember_token"`
	GoogleID          string `json:"google_id"`
	OutletIDDigipos   string `json:"outlet_id_digipos"`
	NamiAgentID       string `json:"nami_agent_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`

	Roles []Role `gorm:"many2many:user_roles"`
}
