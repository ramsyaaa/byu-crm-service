package repository

import (
	"time"

	"gorm.io/gorm"
)

type yaeLeaderboardRepository struct {
	db *gorm.DB
}

func NewYaeLeaderboardRepository(db *gorm.DB) YaeLeaderboardRepository {
	return &yaeLeaderboardRepository{db: db}
}

type VisitHistory struct {
	UserID        int
	AmountDealing float64
}

type LeaderboardItem struct {
	UserID        uint    `json:"user_id"`
	Name          string  `json:"name"`
	AmountDealing float64 `json:"amount_dealing"`
}

func (r *yaeLeaderboardRepository) GetAllLeaderboards(userIDs []int, startDate, endDate time.Time) ([]LeaderboardItem, error) {
	// Kalau tidak ada userID, return kosong
	if len(userIDs) == 0 {
		return []LeaderboardItem{}, nil
	}

	// Query agregasi dengan COALESCE biar NULL jadi 0
	query := `
		SELECT 
			vh.user_id,
			u.name,
			COALESCE(SUM(CAST(JSON_UNQUOTE(JSON_EXTRACT(vh.detail_visit, '$.amount_dealing')) AS DECIMAL(15,2))), 0) AS amount_dealing
		FROM visit_histories vh
		JOIN users u ON u.id = vh.user_id
		WHERE vh.user_id IN ?
		  AND vh.created_at BETWEEN ? AND ?
		GROUP BY vh.user_id, u.name
		ORDER BY amount_dealing DESC
	`

	rows, err := r.db.Raw(query, userIDs, startDate, endDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]LeaderboardItem, 0)
	for rows.Next() {
		var item LeaderboardItem
		if err := rows.Scan(&item.UserID, &item.Name, &item.AmountDealing); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *yaeLeaderboardRepository) GetUserRank(userIDs []int, startDate, endDate time.Time, targetUserID int) (int, int, error) {
	// Kalau tidak ada userID, return kosong
	if len(userIDs) == 0 {
		return 0, 0, nil
	}

	// Query untuk dapatkan ranking user
	query := `
		WITH leaderboard AS (
			SELECT 
				vh.user_id,
				ROW_NUMBER() OVER (
					ORDER BY COALESCE(SUM(CAST(JSON_UNQUOTE(JSON_EXTRACT(vh.detail_visit, '$.amount_dealing')) AS DECIMAL(15,2))), 0) DESC
				) AS rank_position
			FROM visit_histories vh
			WHERE vh.user_id IN ?
			  AND vh.created_at BETWEEN ? AND ?
			GROUP BY vh.user_id
		)
		SELECT rank_position FROM leaderboard WHERE user_id = ?
	`

	var userRank int
	if err := r.db.Raw(query, userIDs, startDate, endDate, targetUserID).Scan(&userRank).Error; err != nil {
		return 0, 0, err
	}

	// Hitung total user yang masuk leaderboard
	var totalUsers int
	if err := r.db.Raw(`
		SELECT COUNT(DISTINCT vh.user_id)
		FROM visit_histories vh
		WHERE vh.user_id IN ?
		  AND vh.created_at BETWEEN ? AND ?
	`, userIDs, startDate, endDate).Scan(&totalUsers).Error; err != nil {
		return 0, 0, err
	}

	return userRank, totalUsers, nil
}
