# Database-Based API Logging System

This document describes the new database-based logging system that replaces the previous file-based logging.

## Overview

The new logging system stores all API requests and responses in a MySQL database table called `api_logs`. This provides better performance, searchability, and analytics capabilities compared to file-based logging.

## Features

### 1. Comprehensive Logging
- **Request Details**: Method, endpoint, IP address, user agent
- **Response Details**: Status code, response time, response body
- **Authentication**: User email extracted from JWT tokens
- **Error Handling**: Detailed error messages and panic recovery
- **Payload Capture**: Request and response payloads (with size limits)

### 2. Asynchronous Processing
- All database writes are performed asynchronously to avoid impacting API performance
- Graceful fallback to console logging if database operations fail

### 3. Performance Optimizations
- Database indexes on frequently queried columns
- Automatic table optimization and analysis
- Configurable payload size limits to prevent storage bloat

### 4. Log Management
- Automatic cleanup of old logs based on retention policies
- Size-based cleanup to maintain database performance
- Periodic optimization of the database table

## Database Schema

The `api_logs` table contains the following fields:

```sql
CREATE TABLE api_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    accessed_at TIMESTAMP NOT NULL,
    endpoint TEXT NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT NOT NULL,
    response_time_ms BIGINT NOT NULL,
    request_payload LONGTEXT,
    response_payload LONGTEXT,
    error_message TEXT,
    auth_user_email VARCHAR(255),
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_accessed_at (accessed_at),
    INDEX idx_method (method),
    INDEX idx_status_code (status_code),
    INDEX idx_auth_user_email (auth_user_email),
    INDEX idx_ip_address (ip_address),
    INDEX idx_created_at (created_at)
);
```

## Configuration

### Environment Variables

You can configure the logging system using these environment variables:

- `LOG_RETENTION_DAYS`: Number of days to keep logs (default: 30)
- `LOG_CLEANUP_INTERVAL_HOURS`: Hours between cleanup runs (default: 24)

### Example .env Configuration

```env
LOG_RETENTION_DAYS=30
LOG_CLEANUP_INTERVAL_HOURS=24
```

## API Endpoints

The system provides several endpoints for viewing and managing logs:

### 1. Get API Logs
```
GET /api-logs?page=1&limit=50&method=GET&status_code=200&user_email=user@example.com&start_date=2024-01-01&end_date=2024-01-31
```

Query Parameters:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50)
- `method`: Filter by HTTP method
- `status_code`: Filter by status code
- `user_email`: Filter by authenticated user email
- `start_date`: Filter from date (YYYY-MM-DD)
- `end_date`: Filter to date (YYYY-MM-DD)

### 2. Get Log Statistics
```
GET /api-logs/stats
```

Returns statistics including:
- Total log count
- Success/error counts
- Average response time
- Date range of logs

### 3. Get Error Logs
```
GET /api-logs/errors?page=1&limit=50
```

Returns only logs with errors (status code >= 400 or error messages).

### 4. Get Slow Requests
```
GET /api-logs/slow?page=1&limit=50&threshold=1000
```

Returns requests slower than the threshold (in milliseconds).

### 5. Get Specific Log
```
GET /api-logs/{id}
```

Returns a specific log entry by ID.

### 6. Manual Cleanup
```
POST /api-logs/cleanup?retention_days=30
```

Manually triggers cleanup of old logs.

## Implementation Details

### JWT User Extraction

The system automatically extracts user email from JWT tokens in the Authorization header:

```go
// Extract email from "Bearer <token>" format
func extractEmailFromJWT(authHeader string) string {
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
    if err != nil {
        return ""
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if email, exists := claims["email"]; exists {
            if emailStr, ok := email.(string); ok {
                return emailStr
            }
        }
    }
    return ""
}
```

### Asynchronous Logging

All database operations are performed in goroutines to avoid blocking the main request flow:

```go
// Log to database asynchronously
go func() {
    if dbErr := db.Create(&logEntry).Error; dbErr != nil {
        // Fallback to console logging if database fails
        log.Printf("Failed to log to database: %v", dbErr)
    }
}()
```

### Error Handling

The system includes comprehensive error handling:
- Database connection failures fall back to console logging
- Panic recovery with stack trace logging
- Graceful handling of malformed JWT tokens

## Migration from File-Based Logging

The old file-based logging system has been replaced, but the log viewer endpoints for file logs are still available for backward compatibility:

- `/log-viewer`: HTML dashboard for file logs
- `/logs`: List available log files
- `/logs/{filename}`: View specific log file content

## Performance Considerations

1. **Payload Size Limits**: Request and response payloads are limited to 10KB to prevent database bloat
2. **Asynchronous Processing**: All database writes are non-blocking
3. **Database Indexes**: Optimized indexes for common query patterns
4. **Automatic Cleanup**: Periodic cleanup prevents unlimited growth
5. **Table Optimization**: Regular OPTIMIZE TABLE operations maintain performance

## Monitoring and Maintenance

The system includes built-in monitoring and maintenance features:

1. **Automatic Index Creation**: Ensures optimal query performance
2. **Periodic Cleanup**: Removes old logs automatically
3. **Table Optimization**: Maintains database performance
4. **Statistics Tracking**: Provides insights into API usage patterns

## Troubleshooting

### Common Issues

1. **Database Connection Errors**: Check database credentials and connectivity
2. **High Storage Usage**: Adjust retention policies or payload size limits
3. **Slow Queries**: Verify indexes are created properly
4. **Missing Logs**: Check for database write permissions and errors in console logs

### Debug Information

Enable debug logging by checking the console output for messages like:
- "Failed to log to database: ..."
- "Cleaned up X old log entries..."
- "Successfully created/verified database indexes..."
