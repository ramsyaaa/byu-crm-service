package response

type LeaderboardData struct {
	UserID        int     `json:"user_id"`
	Name          string  `json:"name"`
	AmountDealing float64 `json:"amount_dealing"`
}
