package response

type BranchResponse struct {
	ID       uint    `json:"id"`
	RegionID int     `json:"region_id"`
	Geojson  *string `json:"geojson" gorm:"type:longtext"`
	Name     string  `json:"name"`
}
