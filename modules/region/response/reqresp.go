package response

type RegionResponse struct {
	ID     uint   `json:"id"`
	AreaID int    `json:"area_id"`
	Name   string `json:"name"`
}
