package repository

import (
	"bytes"
	"byu-crm-service/models"
	"byu-crm-service/modules/notification-one-signal/response"
	"encoding/json"
	"fmt"
	"io"
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

	reqBody := response.NotificationRequest{
		AppID:            OneSignalAppID,
		IncludePlayerIDs: playerID,
		Headings:         map[string]string{"en": requestBody["title"]},
		Contents:         map[string]string{"en": requestBody["description"]},
		URL:              requestBody["callback_url"],
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to send notification, status: %s", string(respBody))
		// return err
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send notification, status: %s", string(respBody))
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
