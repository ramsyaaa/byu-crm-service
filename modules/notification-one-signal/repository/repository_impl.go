package repository

import (
	"bytes"
	"byu-crm-service/models"
	"byu-crm-service/modules/notification-one-signal/response"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type notificationOneSignalRepository struct {
	db *gorm.DB
}

func NewNotificationOneSignalRepository(db *gorm.DB) *notificationOneSignalRepository {
	return &notificationOneSignalRepository{db: db}
}

func (r *notificationOneSignalRepository) SendNotification(requestBody map[string]string, playerID []string) error {
	url := "https://onesignal.com/api/v1/notifications"

	OneSignalAppID := os.Getenv("ONESIGNAL_APP_ID")
	OneSignalAPIKey := os.Getenv("ONESIGNAL_API_KEY")

	callback_url := ""
	if val, ok := requestBody["callback_url"]; ok && val != "" {
		callback_url = val
	}

	reqBody := response.NotificationRequest{
		AppID:            OneSignalAppID,
		IncludePlayerIDs: playerID,
		Headings:         map[string]string{"en": requestBody["title"]},
		Contents:         map[string]string{"en": requestBody["description"]},
		URL:              callback_url,
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+OneSignalAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send notification, status: %s", resp.Status)
	}

	fmt.Println("Notification push successfully!")
	return nil
}

func (r *notificationOneSignalRepository) GetSubscribeNotificationsByUserIDs(userIDs []int) ([]models.SubscribeNotification, error) {
	var subscriptions []models.SubscribeNotification

	if err := r.db.Where("user_id IN ?", userIDs).Find(&subscriptions).Error; err != nil {
		fmt.Println("Error getting subscribe notifications:", err.Error())
		return nil, err
	}

	return subscriptions, nil
}

func (r *notificationOneSignalRepository) CreateSubscribeNotification(dataSubscribe *models.SubscribeNotification) error {
	if err := r.db.Create(dataSubscribe).Error; err != nil {
		fmt.Println("Error creating subscribe notification:", err.Error())
		return err
	}

	return nil
}

func (r *notificationOneSignalRepository) DeleteSubscribeNotification(userID *uint, subscriptionID string) error {
	if err := r.db.Where("user_id = ? AND subscription_id = ?", userID, subscriptionID).
		Delete(&models.SubscribeNotification{}).Error; err != nil {
		return err
	}

	return nil
}
