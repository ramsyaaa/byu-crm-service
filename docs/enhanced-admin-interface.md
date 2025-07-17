# Enhanced Admin Interface Documentation

## Overview

The Enhanced Admin Interface provides a comprehensive dashboard for administrators to monitor Monthly Active Users (MAU) and view API logs. The interface includes authentication, role-based access control, and two main modules with advanced analytics capabilities.

## Features

### 🔐 **Authentication & Security**
- **Admin-only Access**: Restricted to users with `user_type = 'Administrator'`
- **JWT Authentication**: Secure token-based authentication
- **Google OAuth Support**: Alternative login method
- **Session Management**: Automatic token validation and refresh
- **Redirect Protection**: Unauthorized users redirected to login

### 📊 **MAU Dashboard (Module 1)**
- **Monthly Active User Tracking**: Real-time MAU calculations
- **Advanced Filtering**: Date range and user selection filters
- **Interactive Charts**: Daily activity trends and top user visualizations
- **Key Metrics**: Active users, total requests, average per user, peak activity
- **Data Export**: Excel export functionality
- **User Analytics**: Individual user activity timelines

### 📋 **Log Viewer (Module 2)**
- **Enhanced Log Monitoring**: Improved version of existing log viewer
- **Smart Filters**: Advanced filtering capabilities
- **Real-time Analytics**: Live performance metrics
- **Interactive Charts**: Request trends and status distributions
- **Export Functionality**: Excel export of filtered logs

### 🎨 **UI/UX Enhancements**
- **Responsive Design**: Mobile-friendly interface
- **Modern Styling**: Gradient backgrounds and smooth animations
- **Intuitive Navigation**: Tab-based module switching
- **Loading States**: Smooth loading animations
- **Error Handling**: User-friendly error messages

## Architecture

### File Structure
```
static/
├── admin-login.html          # Login page
├── admin-dashboard.html      # Main dashboard container
├── admin-dashboard.js        # Dashboard navigation logic
├── mau-dashboard.html        # MAU module content
├── mau-dashboard.js          # MAU functionality
├── log-viewer-content.html   # Log viewer module content
└── script.js                 # Existing log viewer logic
```

### API Endpoints

#### Authentication
- `POST /api/v1/login` - Regular login
- `GET /api/v1/google/login` - Google OAuth login
- `GET /api/v1/users/profile` - Get user profile (authentication check)

#### MAU Analytics
- `GET /api-logs/mau` - Get MAU data with filters
- `GET /api-logs/users` - Get users list for dropdown
- `GET /api-logs/user-activity` - Get user activity timeline

#### Log Analytics (Existing)
- `GET /api-logs` - Get paginated logs with filters
- `GET /api-logs/stats` - Get log statistics
- `GET /api-logs/chart-data/requests-over-time` - Chart data
- `GET /api-logs/chart-data/status-distribution` - Status distribution

### Routes

#### Admin Routes
- `GET /admin/login` - Login page (public)
- `GET /admin/dashboard` - Main dashboard (protected)
- `GET /admin/` - Redirects to dashboard
- `GET /log-viewer` - Legacy route (redirects to admin dashboard)

## Authentication Flow

### 1. Login Process
1. User accesses `/admin/dashboard`
2. Middleware checks for valid JWT token
3. If no token or invalid, redirect to `/admin/login`
4. User enters credentials or uses Google OAuth
5. Backend validates credentials and user_type
6. If Administrator, JWT token issued and stored
7. Redirect to `/admin/dashboard`

### 2. Authorization Check
```go
// AdminAuthMiddleware checks user_type = 'Administrator'
func AdminAuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // JWT validation
        // Check user_role = "Super-Admin" (maps to Administrator)
        // Store user info in context
        return c.Next()
    }
}
```

## MAU Calculations

### Data Source
- **Table**: `api_logs`
- **Key Fields**: `auth_user_email`, `accessed_at`, `endpoint`
- **Filter**: Only API endpoints (`endpoint LIKE '/api/%'`)

### Metrics Calculated

#### 1. Active Users Count
```sql
SELECT COUNT(DISTINCT auth_user_email) 
FROM api_logs 
WHERE endpoint LIKE '/api/%' 
AND auth_user_email IS NOT NULL 
AND accessed_at BETWEEN start_date AND end_date
```

#### 2. Daily Active Users
```sql
SELECT DATE(accessed_at) as date, COUNT(DISTINCT auth_user_email) as count
FROM api_logs 
WHERE endpoint LIKE '/api/%' 
AND auth_user_email IS NOT NULL
GROUP BY DATE(accessed_at)
ORDER BY date
```

#### 3. Top Active Users
```sql
SELECT auth_user_email, COUNT(*) as request_count, MAX(accessed_at) as last_active
FROM api_logs 
WHERE endpoint LIKE '/api/%' 
AND auth_user_email IS NOT NULL
GROUP BY auth_user_email
ORDER BY request_count DESC
LIMIT 10
```

## Usage Guide

### Accessing the Admin Interface

1. **Navigate to Admin Dashboard**
   ```
   https://your-domain.com/admin/dashboard
   ```

2. **Login with Administrator Account**
   - Use email/password or Google OAuth
   - Only users with `user_type = 'Administrator'` can access

3. **Navigate Between Modules**
   - Click "MAU Dashboard" tab for user analytics
   - Click "Log Viewer" tab for API log monitoring

### Using MAU Dashboard

1. **Set Date Range**
   - Use date pickers or quick filter buttons
   - Default: Current month

2. **Filter by User**
   - Select specific user from dropdown
   - Leave empty for all users

3. **View Metrics**
   - Active Users: Unique users in period
   - Total Requests: API calls made
   - Average per User: Requests divided by users
   - Peak Day: Day with highest activity

4. **Analyze Charts**
   - Daily Active Users: Trend over time
   - Top Users: Most active users bar chart

5. **Export Data**
   - Click "Export Data" for Excel file
   - Includes summary, top users, and daily activity

### Using Log Viewer

1. **Apply Filters**
   - Search: Global text search
   - Method: HTTP method filter
   - Status: Response status codes
   - User: Filter by user email
   - Date Range: Time period
   - Response Time: Performance filters

2. **Quick Filters**
   - Errors Only: Show only failed requests
   - Slow Requests: Show high response times
   - Today Only: Current day only

3. **View Analytics**
   - Real-time metrics cards
   - Interactive charts
   - Detailed log table

4. **Export Logs**
   - Click "Export Excel" for filtered results

## Configuration

### Environment Variables
```env
# JWT Configuration
JWT_SECRET=your-jwt-secret

# Log Retention
LOG_RETENTION_DAYS=30
LOG_CLEANUP_INTERVAL_HOURS=24

# Google OAuth (if used)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### Database Requirements
- Existing `api_logs` table with proper indexes
- User table with `user_type` field
- Proper foreign key relationships

## Security Considerations

### 1. Authentication
- JWT tokens with expiration
- Secure token storage in localStorage
- Automatic token validation

### 2. Authorization
- Role-based access control
- Administrator-only access
- API endpoint protection

### 3. Data Protection
- Sensitive data filtering
- Secure API endpoints
- Input validation and sanitization

## Troubleshooting

### Common Issues

1. **Access Denied**
   - Verify user has `user_type = 'Administrator'`
   - Check JWT token validity
   - Ensure proper role mapping

2. **MAU Data Not Loading**
   - Check `api_logs` table has data
   - Verify date range filters
   - Check database connectivity

3. **Charts Not Displaying**
   - Ensure Chart.js library loaded
   - Check browser console for errors
   - Verify data format from API

4. **Login Issues**
   - Check authentication API endpoints
   - Verify JWT secret configuration
   - Test Google OAuth setup (if used)

### Debug Mode
Enable debug logging by checking browser console for detailed error messages and API response data.

## Future Enhancements

### Planned Features
- Real-time notifications
- Advanced user segmentation
- Custom dashboard widgets
- API rate limiting analytics
- User behavior insights
- Automated reporting

### Performance Optimizations
- Data caching strategies
- Pagination improvements
- Chart rendering optimization
- Background data processing
