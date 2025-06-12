package response

type PermissionResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	GuardName string `json:"guard_name"`
}
