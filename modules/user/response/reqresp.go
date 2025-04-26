package response

type UserResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	Avatar        string   `json:"avatar"`
	Msisdn        string   `json:"msisdn"`
	UserStatus    string   `json:"user_status"`
	UserType      string   `json:"user_type"`
	TerritoryID   uint     `json:"territory_id"`
	TerritoryType string   `json:"territory_type"`
	RoleNames     []string `json:"role_names"`
}
