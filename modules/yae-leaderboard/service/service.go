package service

import (
	"byu-crm-service/modules/yae-leaderboard/response"
	"time"
)

type YaeLeaderboardService interface {
	GetAllLeaderboards(userIDs []int, startDate, endDate time.Time) ([]response.LeaderboardData, error)
}
