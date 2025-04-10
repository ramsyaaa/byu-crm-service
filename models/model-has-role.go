package models

type ModelHasRole struct {
	RoleID    uint   `json:"role_id"`
	ModelID   uint   `json:"model_id"`
	ModelType string `json:"model_type"`
}
