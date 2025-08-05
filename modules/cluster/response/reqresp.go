package response

type ClusterResponse struct {
	ID       uint    `json:"id"`
	BranchID int     `json:"branch_id"`
	Name     string  `json:"name"`
	Geojson  *string `json:"geojson" gorm:"type:longtext"`
}
