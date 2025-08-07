package response

type UserResponse struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Avatar          string   `json:"avatar"`
	Msisdn          string   `json:"msisdn"`
	UserStatus      string   `json:"user_status"`
	UserType        string   `json:"user_type"`
	TerritoryID     uint     `json:"territory_id"`
	TotalPic        *uint    `json:"total_pic"`
	TerritoryType   string   `json:"territory_type"`
	OutletIDDigipos *string  `json:"outlet_id_digipos"`
	NamiAgentID     *string  `json:"nami_agent_id"`
	YaeCode         *string  `json:"yae_code"`
	AreaID          *uint    `json:"area_id"`
	RegionID        *uint    `json:"region_id"`
	BranchID        *uint    `json:"branch_id"`
	ClusterID       *uint    `json:"cluster_id"`
	RoleNames       []string `json:"role_names"`
	Permissions     []string `json:"permissions"`
}
