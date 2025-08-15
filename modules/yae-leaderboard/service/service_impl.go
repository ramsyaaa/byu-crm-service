package service

import (
	"byu-crm-service/modules/yae-leaderboard/repository"
	"byu-crm-service/modules/yae-leaderboard/response"
	"sort"
	"time"
)

type yaeLeaderboardService struct {
	repo repository.YaeLeaderboardRepository
}

func NewYaeLeaderboardService(repo repository.YaeLeaderboardRepository) YaeLeaderboardService {
	return &yaeLeaderboardService{repo: repo}
}

func (s *yaeLeaderboardService) GetAllLeaderboards(userIDs []int, startDate, endDate time.Time) ([]response.LeaderboardData, error) {
	// Step 1: Query leaderboard berdasarkan user ID & tanggal
	visits, err := s.repo.GetAllLeaderboards(userIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Step 2: Urutkan leaderboard
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].AmountDealing > visits[j].AmountDealing
	})

	// Step 3: Convert ke output
	var result []response.LeaderboardData
	for _, v := range visits {
		result = append(result, response.LeaderboardData{
			UserID:        int(v.UserID),
			Name:          v.Name,
			AmountDealing: v.AmountDealing,
		})
	}

	return result, nil
}

func (s *yaeLeaderboardService) GetUserRank(userIDs []int, startDate, endDate time.Time, targetUserID int) (int, int, error) {
	return s.repo.GetUserRank(userIDs, startDate, endDate, targetUserID)
}
