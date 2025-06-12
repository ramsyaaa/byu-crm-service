package models

type RoleHasPermission struct {
	PermissionID uint   `json:"permission_id"`
	RoleID       string `json:"role_id"`
}
