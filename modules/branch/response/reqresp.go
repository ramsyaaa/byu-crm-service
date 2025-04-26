package response

type BranchResponse struct {
	ID       uint   `json:"id"`
	RegionID int    `json:"region_id"`
	Name     string `json:"name"`
}
