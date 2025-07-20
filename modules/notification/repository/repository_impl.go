package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/notification/response"
	"fmt"
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

	if userID > 0 {
		query = query.Where("user_notifications.user_id = ?", userID)
	}

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
	if err != nil {
		return nil, 0, err
	}

	// Modifikasi CallbackUrl
	for i := range user_notifications {
		if user_notifications[i].CallbackUrl != nil {
			parsedURL := *user_notifications[i].CallbackUrl
			separator := "?"
			if strings.Contains(parsedURL, "?") {
				separator = "&"
			}
			newURL := fmt.Sprintf("%s%snotification_id=%d", parsedURL, separator, user_notifications[i].ID)
			user_notifications[i].CallbackUrl = &newURL
		}
	}

	return user_notifications, total, nil

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

func (r *notificationRepository) GetByNotificationId(notificationID uint, userID uint) (*response.NotificationResponse, error) {
	var notification response.NotificationResponse

	query := r.db.Model(&models.UserNotification{}).
		Where("user_notifications.id = ?", notificationID)

	if userID > 0 {
		query = query.Where("user_notifications.user_id = ?", userID)
	}

	err := query.First(&notification).Error
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) MarkNotificationAsRead(notificationID uint, userID uint) error {
	var notification models.UserNotification

	err := r.db.Model(&models.UserNotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		First(&notification).Error
	if err != nil {
		return err
	}

	isRead := uint(1)
	notification.IsRead = &isRead

	return r.db.Save(&notification).Error
}

func (r *notificationRepository) MarkAllNotificationsAsRead(userID int) error {
	var notifications []models.UserNotification

	err := r.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, 0).
		Find(&notifications).Error
	if err != nil {
		return err
	}

	isRead := uint(1)
	for i := range notifications {
		notifications[i].IsRead = &isRead
	}

	return r.db.Save(&notifications).Error
}

func (r *notificationRepository) MarkNotificationAsReadBySubjectID(subjectType string, subjectID uint, userID int) error {
	var notifications []models.UserNotification

	err := r.db.Model(&models.UserNotification{}).
		Where("subject_type = ? AND subject_id = ? AND user_id = ? AND is_read = ?", subjectType, subjectID, userID, 0).
		Find(&notifications).Error
	if err != nil {
		return err
	}

	isRead := uint(1)
	for i := range notifications {
		notifications[i].IsRead = &isRead
	}

	return r.db.Save(&notifications).Error
}
