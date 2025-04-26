package response

type SubdistrictResponse struct {
	ID     uint   `json:"id"`
	CityID int    `json:"city_id"`
	Name   string `json:"name"`
}
