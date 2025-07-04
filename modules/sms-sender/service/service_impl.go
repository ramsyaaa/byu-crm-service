package service

import (
	"bytes"
	userRepo "byu-crm-service/modules/user/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type smsSenderService struct {
	userRepo userRepo.UserRepository
}

func NewSmsSenderService(userRepo userRepo.UserRepository) SmsSenderService {
	return &smsSenderService{userRepo: userRepo}
}

type SmsRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Message     string `json:"message"`
}

func (s *smsSenderService) CreateSms(requestBody map[string]string, rolesName []string, userRole string, territoryID int, userID int) error {
	domainSMS := os.Getenv("DOMAIN_SMS")                                 // Contoh: https://sms.myapp.com
	frontendURL := os.Getenv("FRONTEND_URL")                             // Contoh: https://yae.youthcrm.id
	callbackPath := strings.TrimPrefix(requestBody["callback_url"], "/") // hilangkan '/' jika ada
	frontendURL = fmt.Sprintf("%s/%s", strings.TrimRight(frontendURL, "/"), callbackPath)
	message := requestBody["message"]                                     // Ambil message dari request
	suffix := fmt.Sprintf(" Silakan klik link berikut : %s", frontendURL) // Tambahan di akhir

	if userID != 0 {
		user, err := s.userRepo.FindByID(uint(userID))
		if err != nil {
			return err
		}
		if user.Msisdn != "" && strings.HasPrefix(user.Msisdn, "+62") {
			fullMessage := fmt.Sprintf("%s%s", message, suffix)

			smsPayload := SmsRequest{
				PhoneNumber: user.Msisdn,
				Message:     fullMessage,
			}

			payloadBytes, err := json.Marshal(smsPayload)
			if err != nil {
				return err
			}

			req, err := http.NewRequest("POST", domainSMS+"/api/v1/sms/send", bytes.NewBuffer(payloadBytes))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
		}
	} else {
		filters := map[string]string{
			"search":     "",
			"order_by":   "id",
			"order":      "DESC",
			"start_date": "",
			"end_date":   "",
		}

		users, _, err := s.userRepo.GetAllUsers(0, false, 1, filters, rolesName, false, userRole, territoryID)
		if err != nil {
			return err
		}
		for _, user := range users {
			if user.Msisdn != "" && strings.HasPrefix(user.Msisdn, "+62") {
				fullMessage := fmt.Sprintf("%s%s", message, suffix)

				smsPayload := SmsRequest{
					PhoneNumber: user.Msisdn,
					Message:     fullMessage,
				}

				payloadBytes, err := json.Marshal(smsPayload)
				if err != nil {
					continue
				}

				req, err := http.NewRequest("POST", domainSMS+"/api/v1/sms/send", bytes.NewBuffer(payloadBytes))
				if err != nil {
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				defer resp.Body.Close()
			}
		}

		return nil
	}
	return nil
}
