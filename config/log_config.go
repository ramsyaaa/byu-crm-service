package config

import (
	"log"
	"os"
	"strconv"

	"byu-crm-service/helper"

	"gorm.io/gorm"
)

// LogConfig holds all log-related configuration
type LogConfig struct {
	RetentionDays   int
	CleanupInterval int
	MaxLogSize      int64
	EnableAnalytics bool
}

// LoadLogConfig loads log configuration from environment variables
func LoadLogConfig() *LogConfig {
	config := &LogConfig{
		RetentionDays:   30,      // Default: 30 days
		CleanupInterval: 24,      // Default: 24 hours
		MaxLogSize:      1000000, // Default: 1M logs
		EnableAnalytics: true,
	}

	// Load from environment variables
	if envRetention := os.Getenv("LOG_RETENTION_DAYS"); envRetention != "" {
		if days, err := strconv.Atoi(envRetention); err == nil && days > 0 {
			config.RetentionDays = days
		}
	}

	if envInterval := os.Getenv("LOG_CLEANUP_INTERVAL_HOURS"); envInterval != "" {
		if hours, err := strconv.Atoi(envInterval); err == nil && hours > 0 {
			config.CleanupInterval = hours
		}
	}

	if envMaxSize := os.Getenv("LOG_MAX_SIZE"); envMaxSize != "" {
		if size, err := strconv.ParseInt(envMaxSize, 10, 64); err == nil && size > 0 {
			config.MaxLogSize = size
		}
	}

	if envAnalytics := os.Getenv("LOG_ENABLE_ANALYTICS"); envAnalytics != "" {
		config.EnableAnalytics = envAnalytics == "true"
	}

	return config
}

// InitializeLogManagement sets up log management with the given configuration
func InitializeLogManagement(db *gorm.DB, config *LogConfig) *helper.LogRetentionService {
	// Initialize log management service
	logService := helper.NewLogRetentionService(db)

	// Create database indexes for better performance
	if err := logService.CreateIndexes(); err != nil {
		log.Printf("Failed to create database indexes: %v", err)
	}

	// Start periodic cleanup
	logService.StartPeriodicCleanup(config.RetentionDays, config.CleanupInterval)

	// Log configuration summary
	log.Printf("Log Management Configuration:")
	log.Printf("  - Retention Days: %d", config.RetentionDays)
	log.Printf("  - Cleanup Interval: %d hours", config.CleanupInterval)
	log.Printf("  - Max Log Size: %d entries", config.MaxLogSize)
	log.Printf("  - Analytics Enabled: %t", config.EnableAnalytics)

	return logService
}

// ValidateLogConfig validates the log configuration
func (c *LogConfig) Validate() error {
	if c.RetentionDays <= 0 {
		log.Printf("Warning: Invalid retention days (%d), using default (30)", c.RetentionDays)
		c.RetentionDays = 30
	}

	if c.CleanupInterval <= 0 {
		log.Printf("Warning: Invalid cleanup interval (%d), using default (24)", c.CleanupInterval)
		c.CleanupInterval = 24
	}

	if c.MaxLogSize <= 0 {
		log.Printf("Warning: Invalid max log size (%d), using default (1000000)", c.MaxLogSize)
		c.MaxLogSize = 1000000
	}

	return nil
}

// GetLogConfigSummary returns a summary of current log configuration
func (c *LogConfig) GetLogConfigSummary() map[string]interface{} {
	return map[string]interface{}{
		"retention_days":    c.RetentionDays,
		"cleanup_interval":  c.CleanupInterval,
		"max_log_size":      c.MaxLogSize,
		"analytics_enabled": c.EnableAnalytics,
	}
}
