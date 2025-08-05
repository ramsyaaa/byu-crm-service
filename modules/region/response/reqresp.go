package response

type RegionResponse struct {
	ID      uint    `json:"id"`
	AreaID  int     `json:"area_id"`
	Name    string  `json:"name"`
	Geojson *string `json:"geojson" gorm:"type:longtext"`
}
