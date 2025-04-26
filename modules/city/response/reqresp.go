package response

type CityResponse struct {
	ID        uint   `json:"id"`
	ClusterID int    `json:"cluster_id"`
	Name      string `json:"name"`
}
