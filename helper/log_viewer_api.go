package helper

import (
	"strconv"
	"time"

	"byu-crm-service/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LogViewerHandler handles API log viewing endpoints
type LogViewerHandler struct {
	db *gorm.DB
}

// NewLogViewerHandler creates a new log viewer handler
func NewLogViewerHandler(db *gorm.DB) *LogViewerHandler {
	return &LogViewerHandler{db: db}
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
