package response

type NotificationRequest struct {
	AppID            string            `json:"app_id"`
	IncludePlayerIDs []string          `json:"include_player_ids,omitempty"`
	Headings         map[string]string `json:"headings,omitempty"`
	Contents         map[string]string `json:"contents,omitempty"`
	URL              string            `json:"url,omitempty"` // 👉 tambahkan ini
}
