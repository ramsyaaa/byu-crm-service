package response

type AreaResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RegionResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	AreaID uint   `json:"area_id"`
}

type BranchResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	RegionID uint   `json:"region_id"`
}

type ClusterResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	BranchID uint   `json:"branch_id"`
}

type CityResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ClusterID uint   `json:"cluster_id"`
}
