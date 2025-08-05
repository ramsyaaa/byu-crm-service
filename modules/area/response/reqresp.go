package response

type AreaResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Geojson string `json:"geojson" gorm:"type:longtext"`
}
