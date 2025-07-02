package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/notification/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) GetAllNotifications(filters map[string]string, limit, userID int) ([]response.NotificationResponse, int64, error) {
	var user_notifications []response.NotificationResponse
	var total int64

	query := r.db.Model(&models.UserNotification{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("user_notifications.title LIKE ?", "%"+token+"%"),
			)
		}
	}

	startDate, hasStart := filters["start_date"]
	endDate, hasEnd := filters["end_date"]

	if hasStart && startDate != "" && hasEnd && endDate != "" {
		startDateTime := startDate + " 00:00:00"
		endDateTime := endDate + " 23:59:59"
		query = query.Where("user_notifications.created_at BETWEEN ? AND ?", startDateTime, endDateTime)
	} else if hasStart && startDate != "" {
		startDateTime := startDate + " 00:00:00"
		query = query.Where("user_notifications.created_at >= ?", startDateTime)
	} else if hasEnd && endDate != "" {
		endDateTime := endDate + " 23:59:59"
		query = query.Where("user_notifications.created_at <= ?", endDateTime)
	}

	// Get total count
	query.Count(&total)

	// Ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&user_notifications).Error
	return user_notifications, total, err
}

func (r *notificationRepository) CreateNotification(requestBody map[string]string, userIDs []int) error {
	isRead := uint(0)

	for _, uid := range userIDs {
		u := uint(uid)

		notification := models.UserNotification{
			Title: requestBody["title"],
			Description: func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}(requestBody["description"]),
			CallbackUrl: func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}(requestBody["callback_url"]),
			SubjectType: func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}(requestBody["subject_type"]),
			SubjectID: func(s string) *uint {
				if s == "" {
					return nil
				}
				val, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return nil
				}
				uval := uint(val)
				return &uval
			}(requestBody["subject_id"]),
			UserID: &u,
			IsRead: &isRead,
		}

		if err := r.db.Create(&notification).Error; err != nil {
			return err
		}
	}

	return nil
}
