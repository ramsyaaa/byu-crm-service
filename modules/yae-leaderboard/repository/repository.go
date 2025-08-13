package repository

import "time"

type YaeLeaderboardRepository interface {
	GetAllLeaderboards(userIDs []int, startDate, endDate time.Time) ([]LeaderboardItem, error)
}
