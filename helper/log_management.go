package helper

import (
	"database/sql"
	"log"
	"time"

	"byu-crm-service/models"

	"gorm.io/gorm"
)

// LogRetentionService handles log cleanup and maintenance
type LogRetentionService struct {
	db *gorm.DB
}

// NewLogRetentionService creates a new log retention service
func NewLogRetentionService(db *gorm.DB) *LogRetentionService {
	return &LogRetentionService{db: db}
}

// CleanupOldLogs removes logs older than the specified number of days
func (s *LogRetentionService) CleanupOldLogs(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.ApiLog{})
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Cleaned up %d old log entries older than %d days", result.RowsAffected, retentionDays)
	return nil
}

// CleanupLogsBySize removes oldest logs when total count exceeds maxLogs
func (s *LogRetentionService) CleanupLogsBySize(maxLogs int64) error {
	var count int64
	if err := s.db.Model(&models.ApiLog{}).Count(&count).Error; err != nil {
		return err
	}

	if count <= maxLogs {
		return nil // No cleanup needed
	}

	// Calculate how many logs to delete
	logsToDelete := count - maxLogs

	// Get the oldest logs to delete
	var oldestLogs []models.ApiLog
	if err := s.db.Order("created_at ASC").Limit(int(logsToDelete)).Find(&oldestLogs).Error; err != nil {
		return err
	}

	// Delete the oldest logs
	var ids []uint
	for _, logEntry := range oldestLogs {
		ids = append(ids, logEntry.ID)
	}

	result := s.db.Where("id IN ?", ids).Delete(&models.ApiLog{})
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Cleaned up %d log entries to maintain max size of %d", result.RowsAffected, maxLogs)
	return nil
}

// OptimizeLogTable performs database optimization on the api_logs table
func (s *LogRetentionService) OptimizeLogTable() error {
	// Analyze table for better query performance
	if err := s.db.Exec("ANALYZE TABLE api_logs").Error; err != nil {
		log.Printf("Failed to analyze api_logs table: %v", err)
	}

	// Optimize table to reclaim space after deletions
	if err := s.db.Exec("OPTIMIZE TABLE api_logs").Error; err != nil {
		log.Printf("Failed to optimize api_logs table: %v", err)
		return err
	}

	log.Println("Successfully optimized api_logs table")
	return nil
}

// CreateIndexes creates necessary indexes for better query performance
func (s *LogRetentionService) CreateIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_api_logs_accessed_at ON api_logs(accessed_at)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_method ON api_logs(method)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_status_code ON api_logs(status_code)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_auth_user_email ON api_logs(auth_user_email)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_ip_address ON api_logs(ip_address)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_created_at ON api_logs(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_endpoint_method ON api_logs(endpoint(255), method)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_status_accessed ON api_logs(status_code, accessed_at)",
		// Additional indexes for MAU dashboard performance
		"CREATE INDEX IF NOT EXISTS idx_api_logs_accessed_user_email ON api_logs(accessed_at, auth_user_email)",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_accessed_endpoint ON api_logs(accessed_at, endpoint(255))",
		"CREATE INDEX IF NOT EXISTS idx_api_logs_mau_composite ON api_logs(accessed_at, auth_user_email, endpoint(255))",
		// Note: DATE() function indexes not supported in MariaDB, using accessed_at instead
		"CREATE INDEX IF NOT EXISTS idx_api_logs_accessed_user_composite ON api_logs(accessed_at, auth_user_email, endpoint(100))",
	}

	for _, indexSQL := range indexes {
		if err := s.db.Exec(indexSQL).Error; err != nil {
			log.Printf("Failed to create index: %s, error: %v", indexSQL, err)
		}
	}

	log.Println("Successfully created/verified database indexes for api_logs table")
	return nil
}

// GetLogStats returns statistics about the log table
func (s *LogRetentionService) GetLogStats() (map[string]interface{}, error) {
	return s.GetLogStatsWithFilters("", "", "")
}

// GetLogStatsWithFilters returns statistics about the log table with filters
func (s *LogRetentionService) GetLogStatsWithFilters(apiOnly, startDate, endDate string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Base query
	baseQuery := s.db.Model(&models.ApiLog{})

	// Apply filters
	if apiOnly == "true" {
		baseQuery = baseQuery.Where("endpoint LIKE '/api/%'")
	}
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			baseQuery = baseQuery.Where("accessed_at >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			baseQuery = baseQuery.Where("accessed_at < ?", end.Add(24*time.Hour))
		}
	}

	// Total count
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_logs"] = totalCount

	// Status code counts
	var successCount, clientErrorCount, serverErrorCount int64

	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Where("status_code >= 200 AND status_code < 300").Count(&successCount)

	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Where("status_code >= 400 AND status_code < 500").Count(&clientErrorCount)

	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Where("status_code >= 500").Count(&serverErrorCount)

	stats["success_count"] = successCount
	stats["client_error_count"] = clientErrorCount
	stats["server_error_count"] = serverErrorCount

	// Average response time (handle NULL)
	var avg sql.NullFloat64
	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Select("AVG(response_time_ms)").Scan(&avg)

	if avg.Valid {
		stats["avg_response_time_ms"] = avg.Float64
	} else {
		stats["avg_response_time_ms"] = 0.0
	}

	// Most recent
	var mostRecent time.Time
	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Select("MAX(created_at)").Scan(&mostRecent)
	stats["most_recent_log"] = mostRecent

	// Oldest
	var oldest time.Time
	s.db.Model(&models.ApiLog{}).Scopes(applyFilters(apiOnly, startDate, endDate)).
		Select("MIN(created_at)").Scan(&oldest)
	stats["oldest_log"] = oldest

	return stats, nil
}

// Reusable filter scope
func applyFilters(apiOnly, startDate, endDate string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if apiOnly == "true" {
			db = db.Where("endpoint LIKE '/api/%'")
		}
		if startDate != "" {
			if start, err := time.Parse("2006-01-02", startDate); err == nil {
				db = db.Where("accessed_at >= ?", start)
			}
		}
		if endDate != "" {
			if end, err := time.Parse("2006-01-02", endDate); err == nil {
				db = db.Where("accessed_at < ?", end.Add(24*time.Hour))
			}
		}
		return db
	}
}

// StartPeriodicCleanup starts a background goroutine that periodically cleans up old logs
func (s *LogRetentionService) StartPeriodicCleanup(retentionDays int, cleanupIntervalHours int) {
	go func() {
		ticker := time.NewTicker(time.Duration(cleanupIntervalHours) * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.CleanupOldLogs(retentionDays); err != nil {
					log.Printf("Error during periodic log cleanup: %v", err)
				}

				// Also optimize table periodically
				if err := s.OptimizeLogTable(); err != nil {
					log.Printf("Error during periodic table optimization: %v", err)
				}
			}
		}
	}()

	log.Printf("Started periodic log cleanup: retention=%d days, interval=%d hours", retentionDays, cleanupIntervalHours)
}
