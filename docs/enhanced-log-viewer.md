# Enhanced Log Viewer Dashboard

The Enhanced Log Viewer Dashboard provides a comprehensive, interactive interface for monitoring and analyzing both database-based and file-based logs in your Go Fiber CRM service.

## Features Overview

### 🔄 **Dual Log Source Support**
- **Database Logs**: Real-time API logs stored in MySQL with advanced filtering
- **File Logs**: Legacy file-based logs with backward compatibility

### 🔍 **Advanced Search & Filtering**
- **Global Search**: Search across endpoints, IP addresses, user emails, and error messages
- **Method Filter**: Filter by HTTP methods (GET, POST, PUT, DELETE, etc.)
- **Status Code Filter**: Organized by success (2xx), client error (4xx), and server error (5xx)
- **User Email Filter**: Filter by authenticated user email
- **Date Range Filter**: Select specific date ranges with calendar picker
- **Response Time Filter**: Filter by minimum and maximum response times
- **Quick Filters**: One-click filters for errors, slow requests, and today's logs

### 📊 **Real-time Analytics Dashboard**
- **Total Requests**: Live count of API requests
- **Success Rate**: Percentage of successful requests (2xx status codes)
- **Average Response Time**: Mean response time across all requests
- **Error Rate**: Percentage of server errors (5xx status codes)

### 📈 **Interactive Charts** (Coming Soon)
- **Requests Over Time**: Timeline chart showing request volume
- **Status Code Distribution**: Pie chart of response status codes

### 🎯 **Enhanced User Experience**
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Loading States**: Visual feedback during data loading
- **Empty States**: Clear messaging when no data is found
- **Modal Details**: Click any log entry to view detailed information
- **Pagination**: Advanced pagination with page numbers and navigation
- **Export**: Export filtered results to Excel format

## User Interface Guide

### 1. **Log Source Selection**
At the top of the dashboard, choose between:
- **Database Logs (Real-time)**: Current API logs with full filtering capabilities
- **File Logs (Legacy)**: Historical file-based logs

### 2. **Advanced Filters Panel**
The filters panel includes:

#### Search
- **Global Search Box**: Type to search across endpoints, IPs, users, and errors
- **Real-time Search**: Results update as you type (with 300ms debounce)

#### Method Filter
- Dropdown with all HTTP methods
- Supports multiple selection via Select2 interface

#### Status Code Filter
- Organized by categories (Success, Client Error, Server Error)
- Quick identification of response types

#### User Email Filter
- Search for specific authenticated users
- Supports partial matching

#### Date Range
- **Start Date**: Beginning of date range
- **End Date**: End of date range
- **Calendar Picker**: Easy date selection interface

#### Response Time
- **Min Response Time**: Minimum response time in milliseconds
- **Max Response Time**: Maximum response time in milliseconds

#### Quick Filters
- **Errors Only**: Show only error responses (4xx and 5xx)
- **Slow Requests**: Show requests taking longer than 1 second
- **Today Only**: Show only today's logs

### 3. **Analytics Dashboard**
Four key metrics displayed as cards:
- **Total Requests**: Overall request count
- **Success Rate**: Percentage of successful requests
- **Average Response Time**: Mean response time
- **Error Rate**: Percentage of error responses

### 4. **Data Table**
Enhanced table with:
- **Sortable Columns**: Click headers to sort data
- **Color-coded Status**: Visual status code indicators
- **Truncated Content**: Long URLs and emails are truncated with tooltips
- **Click to View**: Click any row to see detailed information

### 5. **Pagination**
Advanced pagination controls:
- **First/Last Page**: Jump to beginning or end
- **Previous/Next**: Navigate one page at a time
- **Page Numbers**: Direct page navigation
- **Results Counter**: Shows current range and total count
- **Items per Page**: Configurable (25, 50, 100, 200)

### 6. **Action Bar**
- **Refresh Button**: Reload current data
- **Export Button**: Download filtered results as Excel file
- **Results Count**: Shows total matching results

## API Endpoints

The enhanced log viewer uses these API endpoints:

### Database Logs
- `GET /api-logs` - Get paginated logs with filters
- `GET /api-logs/stats` - Get analytics statistics
- `GET /api-logs/errors` - Get error logs only
- `GET /api-logs/slow` - Get slow requests
- `GET /api-logs/:id` - Get specific log entry
- `POST /api-logs/cleanup` - Manual log cleanup

### File Logs (Legacy)
- `GET /logs` - List available log files
- `GET /logs/:filename` - Get content from specific log file

## Filter Parameters

When using database logs, the following query parameters are supported:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `page` | Page number | `?page=2` |
| `limit` | Items per page | `?limit=100` |
| `search` | Global search | `?search=login` |
| `method` | HTTP method | `?method=POST` |
| `status_code` | Status code | `?status_code=500` |
| `user_email` | User email | `?user_email=john@example.com` |
| `start_date` | Start date | `?start_date=2024-01-01` |
| `end_date` | End date | `?end_date=2024-01-31` |
| `min_response_time` | Min response time (ms) | `?min_response_time=1000` |
| `max_response_time` | Max response time (ms) | `?max_response_time=5000` |

## Performance Features

### 1. **Debounced Search**
- Search input is debounced by 300ms to prevent excessive API calls
- Provides smooth user experience while typing

### 2. **Efficient Pagination**
- Server-side pagination reduces data transfer
- Smart page number generation shows relevant pages

### 3. **Optimized Queries**
- Database indexes ensure fast query performance
- Filtered queries reduce unnecessary data processing

### 4. **Asynchronous Loading**
- Non-blocking UI updates
- Loading states provide user feedback

## Browser Compatibility

The enhanced log viewer supports:
- **Chrome**: 80+
- **Firefox**: 75+
- **Safari**: 13+
- **Edge**: 80+

## Dependencies

### Frontend Libraries
- **Tailwind CSS**: Utility-first CSS framework
- **Font Awesome**: Icon library
- **Chart.js**: Charting library
- **Flatpickr**: Date picker
- **Select2**: Enhanced select dropdowns
- **SheetJS**: Excel export functionality
- **jQuery**: DOM manipulation and AJAX

### Backend Dependencies
- **Go Fiber**: Web framework
- **GORM**: ORM for database operations
- **MySQL**: Database storage

## Troubleshooting

### Common Issues

1. **No Data Showing**
   - Check if database logs are enabled
   - Verify database connection
   - Check filter settings

2. **Slow Performance**
   - Reduce date range for large datasets
   - Use more specific filters
   - Check database indexes

3. **Export Not Working**
   - Ensure browser allows downloads
   - Check if data is loaded
   - Try refreshing the page

4. **Filters Not Working**
   - Clear browser cache
   - Check JavaScript console for errors
   - Verify API endpoints are accessible

### Debug Mode

To enable debug mode, open browser developer tools and check the console for:
- API request/response logs
- JavaScript errors
- Network issues

## Future Enhancements

Planned features for future releases:
- **Real-time Updates**: WebSocket-based live log streaming
- **Advanced Charts**: More detailed analytics visualizations
- **Custom Dashboards**: User-configurable dashboard layouts
- **Alert System**: Notifications for error thresholds
- **Log Correlation**: Link related log entries
- **Performance Metrics**: Detailed performance analytics
