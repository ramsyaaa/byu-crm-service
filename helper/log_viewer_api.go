package helper

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"byu-crm-service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Enhanced user data structure with territory information
type EnhancedUserData struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	TerritoryName string `json:"territory_name"`
	TerritoryType string `json:"territory_type"`
}

// Enhanced user activity with full user information
type EnhancedUserActivity struct {
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	TerritoryName string    `json:"territory_name"`
	CallCount     int64     `json:"call_count"`
	LastActivity  time.Time `json:"last_activity"`
}

// LogViewerHandler handles API log viewing endpoints
type LogViewerHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewLogViewerHandler creates a new log viewer handler
func NewLogViewerHandler(db *gorm.DB) *LogViewerHandler {
	return &LogViewerHandler{db: db}
}

// NewLogViewerHandlerWithRedis creates a new log viewer handler with Redis support
func NewLogViewerHandlerWithRedis(db *gorm.DB, redis *redis.Client) *LogViewerHandler {
	return &LogViewerHandler{db: db, redis: redis}
}

// generateCacheKey creates a unique cache key based on query parameters
func (h *LogViewerHandler) generateCacheKey(prefix string, params map[string]string) string {
	// Create a consistent string from parameters by sorting keys
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	// Sort keys to ensure consistent ordering
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	var keyParts []string
	for _, key := range keys {
		keyParts = append(keyParts, fmt.Sprintf("%s=%s", key, params[key]))
	}

	paramStr := fmt.Sprintf("%v", keyParts)

	// Create MD5 hash for shorter, consistent keys
	hash := md5.Sum([]byte(paramStr))
	return fmt.Sprintf("%s:%x", prefix, hash)
}

// getCachedData attempts to retrieve cached data from Redis
func (h *LogViewerHandler) getCachedData(ctx context.Context, cacheKey string) (interface{}, bool) {
	if h.redis == nil {
		return nil, false
	}

	cached, err := h.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("Redis GET error for key %s: %v", cacheKey, err)
		}
		return nil, false
	}

	var cachedData interface{}
	if err := json.Unmarshal([]byte(cached), &cachedData); err != nil {
		log.Printf("Failed to unmarshal cached data for key %s: %v", cacheKey, err)
		return nil, false
	}

	log.Printf("Cache HIT for key: %s", cacheKey)
	return cachedData, true
}

// setCachedData stores data in Redis with expiration
func (h *LogViewerHandler) setCachedData(ctx context.Context, cacheKey string, data interface{}, expiration time.Duration) {
	if h.redis == nil {
		return
	}

	cacheBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data for caching: %v", err)
		return
	}

	if err := h.redis.Set(ctx, cacheKey, cacheBytes, expiration).Err(); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		return
	}

	log.Printf("Cache SET for key: %s (expires in %v)", cacheKey, expiration)
}

// getUsersWithTerritoryInfo retrieves user information with territory names
func (h *LogViewerHandler) getUsersWithTerritoryInfo(emails []string) (map[string]EnhancedUserData, error) {
	if len(emails) == 0 {
		return make(map[string]EnhancedUserData), nil
	}

	var users []struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		TerritoryType string `json:"territory_type"`
		TerritoryID   uint   `json:"territory_id"`
	}

	// Get user data from users table
	if err := h.db.Table("users").
		Select("email, name, territory_type, territory_id").
		Where("email IN ?", emails).
		Scan(&users).Error; err != nil {
		return nil, err
	}

	result := make(map[string]EnhancedUserData)

	// Group users by territory type for efficient querying
	territoryQueries := make(map[string][]uint)
	userTerritoryMap := make(map[string]struct {
		User struct {
			Email         string
			Name          string
			TerritoryType string
		}
		TerritoryID uint
	})

	for _, user := range users {
		territoryQueries[user.TerritoryType] = append(territoryQueries[user.TerritoryType], user.TerritoryID)
		userTerritoryMap[user.Email] = struct {
			User struct {
				Email         string
				Name          string
				TerritoryType string
			}
			TerritoryID uint
		}{
			User: struct {
				Email         string
				Name          string
				TerritoryType string
			}{
				Email:         user.Email,
				Name:          user.Name,
				TerritoryType: user.TerritoryType,
			},
			TerritoryID: user.TerritoryID,
		}
	}

	// Get territory names for each type
	territoryNames := make(map[string]map[uint]string)

	for territoryType, ids := range territoryQueries {
		if len(ids) == 0 {
			continue
		}

		territoryNames[territoryType] = make(map[uint]string)
		var tableName string

		switch territoryType {
		case "App\\Models\\Area":
			tableName = "areas"
		case "App\\Models\\Region":
			tableName = "regions"
		case "App\\Models\\Branch":
			tableName = "branches"
		case "App\\Models\\Cluster":
			tableName = "clusters"
		case "App\\Models\\City":
			tableName = "cities"
		default:
			continue
		}

		var territoryData []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		if err := h.db.Table(tableName).
			Select("id, name").
			Where("id IN ?", ids).
			Scan(&territoryData).Error; err != nil {
			log.Printf("Error fetching territory names from %s: %v", tableName, err)
			continue
		}

		for _, territory := range territoryData {
			territoryNames[territoryType][territory.ID] = territory.Name
		}
	}

	// Build final result
	for email, userData := range userTerritoryMap {
		territoryName := "Not Assigned"
		if names, exists := territoryNames[userData.User.TerritoryType]; exists {
			if name, found := names[userData.TerritoryID]; found {
				territoryName = name
			}
		}

		result[email] = EnhancedUserData{
			Email:         userData.User.Email,
			Name:          userData.User.Name,
			TerritoryName: territoryName,
			TerritoryType: userData.User.TerritoryType,
		}
	}

	// Add entries for emails not found in users table
	for _, email := range emails {
		if _, exists := result[email]; !exists {
			result[email] = EnhancedUserData{
				Email:         email,
				Name:          email, // Fallback to email as name
				TerritoryName: "Not Assigned",
				TerritoryType: "",
			}
		}
	}

	return result, nil
}

// GetApiLogs returns paginated API logs with filtering options
func (h *LogViewerHandler) GetApiLogs(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	method := c.Query("method")
	statusCode := c.Query("status_code")
	userEmail := c.Query("user_email")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	minResponseTime := c.Query("min_response_time")
	maxResponseTime := c.Query("max_response_time")
	apiOnly := c.Query("api_only")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := h.db.Model(&models.ApiLog{})

	// Filter only API endpoints if requested
	if apiOnly == "true" {
		query = query.Where("endpoint LIKE '/api/%'")
	}

	// Apply search filter (searches across endpoint, IP, user email)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"endpoint LIKE ? OR ip_address LIKE ? OR auth_user_email LIKE ? OR error_message LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Apply filters
	if method != "" {
		query = query.Where("method = ?", method)
	}

	if statusCode != "" {
		if code, err := strconv.Atoi(statusCode); err == nil {
			query = query.Where("status_code = ?", code)
		}
	}

	if userEmail != "" {
		query = query.Where("auth_user_email LIKE ?", "%"+userEmail+"%")
	}

	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("accessed_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 24 hours to include the entire end date
			endDateTime := end.Add(24 * time.Hour)
			query = query.Where("accessed_at < ?", endDateTime)
		}
	}

	// Apply response time filters
	if minResponseTime != "" {
		if minTime, err := strconv.Atoi(minResponseTime); err == nil {
			query = query.Where("response_time_ms >= ?", minTime)
		}
	}

	if maxResponseTime != "" {
		if maxTime, err := strconv.Atoi(maxResponseTime); err == nil {
			query = query.Where("response_time_ms <= ?", maxTime)
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get logs with pagination
	var logs []models.ApiLog
	result := query.Order("accessed_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch logs",
			"error":   result.Error.Error(),
		})
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"logs": logs,
			"pagination": fiber.Map{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  total,
				"limit":        limit,
			},
		},
	})
}

// GetLogStats returns statistics about API logs
func (h *LogViewerHandler) GetLogStats(c *fiber.Ctx) error {
	apiOnly := c.Query("api_only")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	logService := NewLogRetentionService(h.db)
	stats, err := logService.GetLogStatsWithFilters(apiOnly, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch log statistics",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

// GetLogById returns a specific log entry by ID
func (h *LogViewerHandler) GetLogById(c *fiber.Ctx) error {
	id := c.Params("id")

	var log models.ApiLog
	result := h.db.First(&log, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Log entry not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch log entry",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   log,
	})
}

// CleanupLogs manually triggers log cleanup
func (h *LogViewerHandler) CleanupLogs(c *fiber.Ctx) error {
	retentionDays, _ := strconv.Atoi(c.Query("retention_days", "30"))

	logService := NewLogRetentionService(h.db)
	err := logService.CleanupOldLogs(retentionDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to cleanup logs",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Log cleanup completed successfully",
	})
}

// GetErrorLogs returns logs with errors only
func (h *LogViewerHandler) GetErrorLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	offset := (page - 1) * limit

	// Query for logs with errors (status code >= 400 or error_message is not null)
	query := h.db.Model(&models.ApiLog{}).Where("status_code >= 400 OR error_message IS NOT NULL")

	var total int64
	query.Count(&total)

	var logs []models.ApiLog
	result := query.Order("accessed_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch error logs",
			"error":   result.Error.Error(),
		})
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"logs": logs,
			"pagination": fiber.Map{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  total,
				"limit":        limit,
			},
		},
	})
}

// GetSlowRequests returns logs with slow response times
func (h *LogViewerHandler) GetSlowRequests(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	threshold, _ := strconv.Atoi(c.Query("threshold", "1000")) // Default 1 second

	offset := (page - 1) * limit

	query := h.db.Model(&models.ApiLog{}).Where("response_time_ms > ?", threshold)

	var total int64
	query.Count(&total)

	var logs []models.ApiLog
	result := query.Order("response_time_ms DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch slow requests",
			"error":   result.Error.Error(),
		})
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"logs": logs,
			"pagination": fiber.Map{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  total,
				"limit":        limit,
			},
		},
	})
}

// GetRequestsOverTime returns data for requests over time chart
func (h *LogViewerHandler) GetRequestsOverTime(c *fiber.Ctx) error {
	apiOnly := c.Query("api_only")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Build query
	query := h.db.Model(&models.ApiLog{})

	// Apply API only filter
	if apiOnly == "true" {
		query = query.Where("endpoint LIKE '/api/%'")
	}

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("accessed_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := end.Add(24 * time.Hour)
			query = query.Where("accessed_at < ?", endDateTime)
		}
	}

	// Group by hour for requests over time
	var results []struct {
		Hour  string `json:"hour"`
		Count int64  `json:"count"`
	}

	err := query.Select("DATE_FORMAT(accessed_at, '%H:00') as hour, COUNT(*) as count").
		Group("DATE_FORMAT(accessed_at, '%H:00')").
		Order("hour").
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch requests over time data",
			"error":   err.Error(),
		})
	}

	// Prepare chart data
	labels := make([]string, 0)
	values := make([]int64, 0)

	for _, result := range results {
		labels = append(labels, result.Hour)
		values = append(values, result.Count)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"labels": labels,
			"values": values,
		},
	})
}

// GetStatusDistribution returns data for status code distribution chart
func (h *LogViewerHandler) GetStatusDistribution(c *fiber.Ctx) error {
	apiOnly := c.Query("api_only")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Build query
	query := h.db.Model(&models.ApiLog{})

	// Apply API only filter
	if apiOnly == "true" {
		query = query.Where("endpoint LIKE '/api/%'")
	}

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("accessed_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := end.Add(24 * time.Hour)
			query = query.Where("accessed_at < ?", endDateTime)
		}
	}

	// Count by status code ranges
	var successCount, clientErrorCount, serverErrorCount int64
	query.Where("status_code >= 200 AND status_code < 300").Count(&successCount)
	query.Where("status_code >= 400 AND status_code < 500").Count(&clientErrorCount)
	query.Where("status_code >= 500").Count(&serverErrorCount)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"success_count":      successCount,
			"client_error_count": clientErrorCount,
			"server_error_count": serverErrorCount,
		},
	})
}

// GetMAUStats returns Monthly Active Users statistics
func (h *LogViewerHandler) GetMAUStats(c *fiber.Ctx) error {
	// Parse query parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	apiOnly := c.Query("api_only", "true")              // Default to API only for MAU
	businessHours := c.Query("business_hours", "false") // Default to all hours
	uniqueOnly := c.Query("unique_only", "true")        // Default to unique calls only

	// Default to today only if no dates provided (for better performance)
	now := time.Now()
	if startDate == "" {
		startDate = now.Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	// Create cache key based on all parameters
	cacheParams := map[string]string{
		"start_date":     startDate,
		"end_date":       endDate,
		"api_only":       apiOnly,
		"business_hours": businessHours,
		"unique_only":    uniqueOnly,
	}
	cacheKey := h.generateCacheKey("mau_stats", cacheParams)

	// Try to get from cache first
	ctx := context.Background()
	if cachedData, found := h.getCachedData(ctx, cacheKey); found {
		// Add cached flag to the existing response structure
		if dataMap, ok := cachedData.(map[string]interface{}); ok {
			dataMap["cached"] = true
			return c.JSON(fiber.Map{
				"status": "success",
				"data":   dataMap,
			})
		}
		// Fallback if cached data structure is unexpected
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"cached": true,
			},
		})
	}

	log.Printf("Cache MISS for key: %s - executing database query", cacheKey)

	// Parse dates
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid end_date format. Use YYYY-MM-DD",
		})
	}
	endDateTime = endDateTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // End of day

	// Build optimized base query with proper indexing (exclude admin account)
	baseWhere := "accessed_at >= ? AND accessed_at <= ? AND auth_user_email IS NOT NULL AND auth_user_email != '' AND auth_user_email != 'super@admin.com'"
	var baseArgs []interface{}
	baseArgs = append(baseArgs, startDateTime, endDateTime)

	if apiOnly == "true" {
		baseWhere += " AND endpoint LIKE ?"
		baseArgs = append(baseArgs, "/api/%")
	}

	// Add business hours filtering (8 AM - 7 PM)
	if businessHours == "true" {
		baseWhere += " AND HOUR(accessed_at) >= 8 AND HOUR(accessed_at) <= 19"
	}

	// Use a single optimized query to get both current MAU and total calls
	type MAUResult struct {
		CurrentMAU    int64 `json:"current_mau"`
		TotalAPICalls int64 `json:"total_api_calls"`
	}

	var result MAUResult
	var sqlQuery string

	if uniqueOnly == "true" {
		// Count unique API calls per user per endpoint per day
		sqlQuery = `
			SELECT
				COUNT(DISTINCT auth_user_email) as current_mau,
				COUNT(DISTINCT CONCAT(auth_user_email, '|', endpoint, '|', DATE(accessed_at))) as total_api_calls
			FROM api_logs
			WHERE ` + baseWhere
	} else {
		// Count all API calls
		sqlQuery = `
			SELECT
				COUNT(DISTINCT auth_user_email) as current_mau,
				COUNT(*) as total_api_calls
			FROM api_logs
			WHERE ` + baseWhere
	}

	if err := h.db.Raw(sqlQuery, baseArgs...).Scan(&result).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch MAU statistics",
			"error":   err.Error(),
		})
	}

	// Get previous period MAU for comparison (only if not today-only)
	var previousMAU int64
	var growthPercentage float64

	// Calculate previous period based on the date range
	duration := endDateTime.Sub(startDateTime)
	prevEndDateTime := startDateTime.Add(-time.Second)
	prevStartDateTime := prevEndDateTime.Add(-duration)

	prevWhere := "accessed_at >= ? AND accessed_at <= ? AND auth_user_email IS NOT NULL AND auth_user_email != '' AND auth_user_email != 'super@admin.com'"
	var prevArgs []interface{}
	prevArgs = append(prevArgs, prevStartDateTime, prevEndDateTime)

	if apiOnly == "true" {
		prevWhere += " AND endpoint LIKE ?"
		prevArgs = append(prevArgs, "/api/%")
	}

	// Add business hours filtering for previous period
	if businessHours == "true" {
		prevWhere += " AND HOUR(accessed_at) >= 8 AND HOUR(accessed_at) <= 19"
	}

	err = h.db.Raw(`
		SELECT COUNT(DISTINCT auth_user_email) as previous_mau
		FROM api_logs
		WHERE `+prevWhere, prevArgs...).Scan(&previousMAU).Error

	if err == nil {
		// Calculate growth percentage
		if previousMAU > 0 {
			growthPercentage = ((float64(result.CurrentMAU) - float64(previousMAU)) / float64(previousMAU)) * 100
		} else if result.CurrentMAU > 0 {
			growthPercentage = 100 // 100% growth from 0
		}
	}

	// Prepare response data for caching
	responseData := fiber.Map{
		"current_mau":       result.CurrentMAU,
		"previous_mau":      previousMAU,
		"growth_percentage": growthPercentage,
		"total_api_calls":   result.TotalAPICalls,
		"period_start":      startDate,
		"period_end":        endDate,
		"business_hours":    businessHours == "true",
		"unique_only":       uniqueOnly == "true",
		"api_only":          apiOnly == "true",
	}

	// Cache the results for 12 hours
	cacheExpiration := 12 * time.Hour
	h.setCachedData(ctx, cacheKey, responseData, cacheExpiration)

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   responseData,
	})
}

// GetUserActivityData returns user activity data for MAU dashboard
func (h *LogViewerHandler) GetUserActivityData(c *fiber.Ctx) error {
	// Parse query parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	apiOnly := c.Query("api_only", "true")
	businessHours := c.Query("business_hours", "false") // Default to all hours
	uniqueOnly := c.Query("unique_only", "true")        // Default to unique calls only
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Default to today only if no dates provided (for better performance)
	now := time.Now()
	if startDate == "" {
		startDate = now.Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	// Create cache key based on all parameters
	cacheParams := map[string]string{
		"start_date":     startDate,
		"end_date":       endDate,
		"api_only":       apiOnly,
		"business_hours": businessHours,
		"unique_only":    uniqueOnly,
		"limit":          strconv.Itoa(limit),
	}
	cacheKey := h.generateCacheKey("user_activity", cacheParams)

	// Try to get from cache first
	ctx := context.Background()
	if cachedData, found := h.getCachedData(ctx, cacheKey); found {
		// Add cached flag to the existing response structure
		if dataMap, ok := cachedData.(map[string]interface{}); ok {
			dataMap["cached"] = true
			return c.JSON(fiber.Map{
				"status": "success",
				"data":   dataMap,
			})
		}
		// Fallback if cached data structure is unexpected
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"cached": true,
			},
		})
	}

	log.Printf("Cache MISS for key: %s - executing database query", cacheKey)

	// Parse dates
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid end_date format. Use YYYY-MM-DD",
		})
	}
	endDateTime = endDateTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Build optimized query for top active users (exclude admin account)
	baseWhere := "accessed_at >= ? AND accessed_at <= ? AND auth_user_email IS NOT NULL AND auth_user_email != '' AND auth_user_email != 'super@admin.com'"
	var baseArgs []interface{}
	baseArgs = append(baseArgs, startDateTime, endDateTime)

	if apiOnly == "true" {
		baseWhere += " AND endpoint LIKE ?"
		baseArgs = append(baseArgs, "/api/%")
	}

	// Add business hours filtering (8 AM - 7 PM)
	if businessHours == "true" {
		baseWhere += " AND HOUR(accessed_at) >= 8 AND HOUR(accessed_at) <= 19"
	}

	type UserActivity struct {
		AuthUserEmail string    `json:"auth_user_email"`
		CallCount     int64     `json:"call_count"`
		LastActivity  time.Time `json:"last_activity"`
	}

	var userActivities []UserActivity
	var sqlQuery string

	if uniqueOnly == "true" {
		// Count unique API calls per user per endpoint per day
		sqlQuery = `
			SELECT
				auth_user_email,
				COUNT(DISTINCT CONCAT(endpoint, '|', DATE(accessed_at))) as call_count,
				MAX(accessed_at) as last_activity
			FROM api_logs
			WHERE ` + baseWhere + `
			GROUP BY auth_user_email
			ORDER BY call_count DESC
			LIMIT ?`
	} else {
		// Count all API calls
		sqlQuery = `
			SELECT
				auth_user_email,
				COUNT(*) as call_count,
				MAX(accessed_at) as last_activity
			FROM api_logs
			WHERE ` + baseWhere + `
			GROUP BY auth_user_email
			ORDER BY call_count DESC
			LIMIT ?`
	}

	if err := h.db.Raw(sqlQuery, append(baseArgs, limit)...).Scan(&userActivities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch user activity data",
			"error":   err.Error(),
		})
	}

	// Extract user emails for enhanced data lookup
	var userEmails []string
	for _, activity := range userActivities {
		userEmails = append(userEmails, activity.AuthUserEmail)
	}

	// Get enhanced user data with territory information
	enhancedUsers, err := h.getUsersWithTerritoryInfo(userEmails)
	if err != nil {
		log.Printf("Error fetching enhanced user data: %v", err)
		// Fallback to original data
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"top_users":    userActivities,
				"period_start": startDate,
				"period_end":   endDate,
			},
		})
	}

	// Build enhanced user activity data
	var enhancedActivities []EnhancedUserActivity
	for _, activity := range userActivities {
		enhanced := EnhancedUserActivity{
			Email:         activity.AuthUserEmail,
			Name:          activity.AuthUserEmail, // Default fallback
			TerritoryName: "Not Assigned",         // Default fallback
			CallCount:     activity.CallCount,
			LastActivity:  activity.LastActivity,
		}

		if userData, exists := enhancedUsers[activity.AuthUserEmail]; exists {
			enhanced.Name = userData.Name
			enhanced.TerritoryName = userData.TerritoryName
		}

		enhancedActivities = append(enhancedActivities, enhanced)
	}

	// Prepare response data for caching
	responseData := fiber.Map{
		"top_users":    enhancedActivities,
		"period_start": startDate,
		"period_end":   endDate,
	}

	// Cache the results for 12 hours
	cacheExpiration := 12 * time.Hour
	h.setCachedData(ctx, cacheKey, responseData, cacheExpiration)

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   responseData,
	})
}

// GetDailyActiveUsers returns daily active users trend data with Redis caching
func (h *LogViewerHandler) GetDailyActiveUsers(c *fiber.Ctx) error {
	// Parse query parameters
	days, _ := strconv.Atoi(c.Query("days", "30")) // Default to last 30 days
	apiOnly := c.Query("api_only", "true")
	businessHours := c.Query("business_hours", "false")
	uniqueOnly := c.Query("unique_only", "true")
	userFilter := c.Query("user_filter", "") // Optional user filter

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Create cache key based on all parameters
	cacheParams := map[string]string{
		"days":           strconv.Itoa(days),
		"api_only":       apiOnly,
		"business_hours": businessHours,
		"unique_only":    uniqueOnly,
		"user_filter":    userFilter,
		"start_date":     startDate.Format("2006-01-02"),
		"end_date":       endDate.Format("2006-01-02"),
	}
	cacheKey := h.generateCacheKey("daily_active_users", cacheParams)

	// Try to get from cache first
	ctx := context.Background()
	if cachedData, found := h.getCachedData(ctx, cacheKey); found {
		// Add cached flag to the existing response structure
		if dataMap, ok := cachedData.(map[string]interface{}); ok {
			dataMap["cached"] = true
			return c.JSON(fiber.Map{
				"status": "success",
				"data":   dataMap,
			})
		}
		// Fallback if cached data structure is unexpected
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"cached": true,
			},
		})
	}

	log.Printf("Cache MISS for key: %s - executing database query", cacheKey)

	// Build optimized query for daily active users (exclude admin account)
	baseWhere := "accessed_at >= ? AND accessed_at <= ? AND auth_user_email IS NOT NULL AND auth_user_email != '' AND auth_user_email != 'super@admin.com'"
	var baseArgs []interface{}
	baseArgs = append(baseArgs, startDate, endDate)

	if apiOnly == "true" {
		baseWhere += " AND endpoint LIKE ?"
		baseArgs = append(baseArgs, "/api/%")
	}

	// Add business hours filtering (8 AM - 7 PM)
	if businessHours == "true" {
		baseWhere += " AND HOUR(accessed_at) >= 8 AND HOUR(accessed_at) <= 19"
	}

	// Add user filter if specified
	if userFilter != "" {
		baseWhere += " AND auth_user_email = ?"
		baseArgs = append(baseArgs, userFilter)
	}

	type DailyActivity struct {
		Date        string `json:"date"`
		ActiveUsers int64  `json:"active_users"`
	}

	var dailyActivities []DailyActivity
	var sqlQuery string
	var queryArgs []interface{}

	// Choose query based on unique_only parameter
	if uniqueOnly == "true" {
		// Count unique users per day (not user-endpoint combinations)
		sqlQuery = `
			SELECT
				DATE(accessed_at) as date,
				COUNT(DISTINCT auth_user_email) as active_users
			FROM api_logs
			WHERE ` + baseWhere + `
			GROUP BY DATE(accessed_at)
			ORDER BY date ASC`
		queryArgs = baseArgs
	} else {
		// Use the optimized CTE approach for regular counting
		additionalFilters := ""
		if apiOnly == "true" {
			additionalFilters += "AND al.endpoint LIKE '/api/%'"
		}
		if businessHours == "true" {
			additionalFilters += " AND HOUR(al.accessed_at) >= 8 AND HOUR(al.accessed_at) <= 19"
		}
		if userFilter != "" {
			additionalFilters += " AND al.auth_user_email = '" + userFilter + "'"
		}

		sqlQuery = `
			WITH RECURSIVE date_range AS (
				SELECT DATE(?) as date
				UNION ALL
				SELECT DATE_ADD(date, INTERVAL 1 DAY)
				FROM date_range
				WHERE date < DATE(?)
			)
			SELECT
				dr.date,
				COALESCE(COUNT(DISTINCT al.auth_user_email), 0) as active_users
			FROM date_range dr
			LEFT JOIN api_logs al ON DATE(al.accessed_at) = dr.date
				AND al.accessed_at >= ?
				AND al.accessed_at <= ?
				AND al.auth_user_email IS NOT NULL
				AND al.auth_user_email != ''
				AND al.auth_user_email != 'super@admin.com'
				` + additionalFilters + `
			GROUP BY dr.date
			ORDER BY dr.date ASC`
		queryArgs = []interface{}{startDate, endDate, startDate, endDate}
	}

	// Execute the query
	if err := h.db.Raw(sqlQuery, queryArgs...).Scan(&dailyActivities).Error; err != nil {
		// Fallback to simpler query if CTE is not supported
		fallbackQuery := `
			SELECT
				DATE(accessed_at) as date,
				COUNT(DISTINCT auth_user_email) as active_users
			FROM api_logs
			WHERE ` + baseWhere + `
			GROUP BY DATE(accessed_at)
			ORDER BY date ASC`

		if err := h.db.Raw(fallbackQuery, baseArgs...).Scan(&dailyActivities).Error; err != nil {
			log.Printf("Database query failed for daily active users: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch daily activity data",
				"error":   err.Error(),
			})
		}
	}

	// Prepare response data for caching
	responseData := fiber.Map{
		"daily_activities": dailyActivities,
		"period_days":      days,
		"start_date":       startDate.Format("2006-01-02"),
		"end_date":         endDate.Format("2006-01-02"),
		"business_hours":   businessHours == "true",
		"unique_only":      uniqueOnly == "true",
		"api_only":         apiOnly == "true",
	}

	// Cache the results for 12 hours (43,200 seconds)
	cacheExpiration := 12 * time.Hour
	h.setCachedData(ctx, cacheKey, responseData, cacheExpiration)

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   responseData,
	})
}

// GetActiveUsersList returns list of active users for dropdown filter
func (h *LogViewerHandler) GetActiveUsersList(c *fiber.Ctx) error {
	// Parse query parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	apiOnly := c.Query("api_only", "true")
	businessHours := c.Query("business_hours", "false") // Default to all hours

	// Default to today only if no dates provided (for better performance)
	now := time.Now()
	if startDate == "" {
		startDate = now.Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	// Parse dates
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid end_date format. Use YYYY-MM-DD",
		})
	}
	endDateTime = endDateTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Create cache key based on all parameters
	cacheParams := map[string]string{
		"start_date":     startDate,
		"end_date":       endDate,
		"api_only":       apiOnly,
		"business_hours": businessHours,
	}
	cacheKey := h.generateCacheKey("active_users_list", cacheParams)

	// Try to get from cache first
	ctx := context.Background()
	if cachedData, found := h.getCachedData(ctx, cacheKey); found {
		// Add cached flag to the existing response structure
		if dataMap, ok := cachedData.(map[string]interface{}); ok {
			dataMap["cached"] = true
			return c.JSON(fiber.Map{
				"status": "success",
				"data":   dataMap,
			})
		}
		// Fallback if cached data structure is unexpected
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"cached": true,
			},
		})
	}

	log.Printf("Cache MISS for key: %s - executing database query", cacheKey)

	// Build optimized query for distinct user emails (exclude admin account)
	baseWhere := "accessed_at >= ? AND accessed_at <= ? AND auth_user_email IS NOT NULL AND auth_user_email != '' AND auth_user_email != 'super@admin.com'"
	var baseArgs []interface{}
	baseArgs = append(baseArgs, startDateTime, endDateTime)

	if apiOnly == "true" {
		baseWhere += " AND endpoint LIKE ?"
		baseArgs = append(baseArgs, "/api/%")
	}

	// Add business hours filtering (8 AM - 7 PM)
	if businessHours == "true" {
		baseWhere += " AND HOUR(accessed_at) >= 8 AND HOUR(accessed_at) <= 19"
	}

	var userEmails []string
	if err := h.db.Raw(`
		SELECT DISTINCT auth_user_email
		FROM api_logs
		WHERE `+baseWhere+`
		ORDER BY auth_user_email ASC`, baseArgs...).Pluck("auth_user_email", &userEmails).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch user list",
			"error":   err.Error(),
		})
	}

	// Get enhanced user data with territory information
	enhancedUsers, err := h.getUsersWithTerritoryInfo(userEmails)
	if err != nil {
		log.Printf("Error fetching enhanced user data: %v", err)
		// Fallback to simple email list
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"users":        userEmails,
				"total_users":  len(userEmails),
				"period_start": startDate,
				"period_end":   endDate,
			},
		})
	}

	// Convert map to slice for consistent ordering
	var enhancedUserList []EnhancedUserData
	for _, email := range userEmails {
		if userData, exists := enhancedUsers[email]; exists {
			enhancedUserList = append(enhancedUserList, userData)
		}
	}

	// Prepare response data for caching
	responseData := fiber.Map{
		"users":        enhancedUserList,
		"total_users":  len(enhancedUserList),
		"period_start": startDate,
		"period_end":   endDate,
	}

	// Cache the results for 12 hours
	cacheExpiration := 12 * time.Hour
	h.setCachedData(ctx, cacheKey, responseData, cacheExpiration)

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   responseData,
	})
}
