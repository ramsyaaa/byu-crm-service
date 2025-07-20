package service

import (
	"bytes"
	userRepo "byu-crm-service/modules/user/repository"
	"encoding/json"
	"fmt"
	"io"
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

func (s *smsSenderService) AssignSmsToUsers(requestBody map[string]string, userIDs []int) error {
	domainSMS := os.Getenv("DOMAIN_SMS") // Contoh: https://sms.myapp.com
	usingLink := true
	frontendURL := os.Getenv("FRONTEND_URL") // Contoh: https://yae.youthcrm.id
	callbackPath := ""
	if requestBody["callback_url"] == "" {
		usingLink = false
	} else {
		frontendURL = os.Getenv("FRONTEND_URL")                             // Contoh: https://yae.youthcrm.id
		callbackPath = strings.TrimPrefix(requestBody["callback_url"], "/") // hilangkan '/' jika ada
		frontendURL = fmt.Sprintf("%s/%s", strings.TrimRight(frontendURL, "/"), callbackPath)
	}

	message := requestBody["message"] // Ambil message dari request

	suffix := fmt.Sprintf(" Silakan klik link berikut : %s", frontendURL) // Tambahan di akhir

	// Convert []int to []uint
	userIDsUint := make([]uint, len(userIDs))
	for i, id := range userIDs {
		userIDsUint[i] = uint(id)
	}
	users, err := s.userRepo.GetUserByIDs(userIDsUint)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Msisdn == "" {
			continue
		}

		fmt.Println("User MSISDN:", user.Msisdn)

		msisdn := normalizeMsisdn(user.Msisdn)
		if msisdn == "" {
			continue // skip kalau gagal normalisasi
		}

		fmt.Println("Normalized MSISDN:", msisdn)

		fullMessage := ""
		if usingLink {
			fullMessage = fmt.Sprintf("%s%s", message, suffix)
		} else {
			fullMessage = message
		}

		smsPayload := SmsRequest{
			PhoneNumber: msisdn,
			Message:     fullMessage,
		}

		payloadBytes, err := json.Marshal(smsPayload)
		if err != nil {
			continue
		}

		fmt.Println("domaim", domainSMS)
		req, err := http.NewRequest("POST", domainSMS+"/api/v1/sms/send", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		fmt.Println("RESPONSE", resp)
		defer resp.Body.Close()
	}

	return nil
}

func (s *smsSenderService) CreateSms(requestBody map[string]string, rolesName []string, userRole string, territoryID int, userID int) error {
	domainSMS := os.Getenv("DOMAIN_SMS") // Contoh: https://sms.myapp.com
	usingLink := true
	frontendURL := os.Getenv("FRONTEND_URL") // Contoh: https://yae.youthcrm.id
	callbackPath := ""
	if requestBody["callback_url"] == "" {
		usingLink = false
	} else {
		frontendURL = os.Getenv("FRONTEND_URL")                             // Contoh: https://yae.youthcrm.id
		callbackPath = strings.TrimPrefix(requestBody["callback_url"], "/") // hilangkan '/' jika ada
		frontendURL = fmt.Sprintf("%s/%s", strings.TrimRight(frontendURL, "/"), callbackPath)
	}

	message := requestBody["message"]                                     // Ambil message dari request
	suffix := fmt.Sprintf(" Silakan klik link berikut : %s", frontendURL) // Tambahan di akhir

	if userID != 0 {
		user, err := s.userRepo.FindByID(uint(userID))
		if err != nil {
			return err
		}
		if user.Msisdn != "" && strings.HasPrefix(user.Msisdn, "+62") {
			fullMessage := ""
			if usingLink {
				fullMessage = fmt.Sprintf("%s%s", message, suffix)
			} else {
				fullMessage = fmt.Sprintf("%s", message)
			}

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
			if user.Msisdn == "" {
				continue
			}

			msisdn := normalizeMsisdn(user.Msisdn)
			if msisdn == "" {
				continue // skip kalau gagal normalisasi
			}

			fullMessage := ""
			if usingLink {
				fullMessage = fmt.Sprintf("%s%s", message, suffix)
			} else {
				fullMessage = message
			}

			smsPayload := SmsRequest{
				PhoneNumber: msisdn,
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
				fmt.Println("Error sending request:", err)
				continue
			}
			fmt.Println("RESPONSE", resp)
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
			} else {
				fmt.Println("Error response body:", string(body))
			}
		}

		return nil
	}
	return nil
}

func normalizeMsisdn(msisdn string) string {
	msisdn = strings.TrimSpace(msisdn)
	msisdn = strings.ReplaceAll(msisdn, " ", "")
	msisdn = strings.ReplaceAll(msisdn, "-", "")

	if strings.HasPrefix(msisdn, "+62") {
		return msisdn
	} else if strings.HasPrefix(msisdn, "62") {
		return "+" + msisdn
	} else if strings.HasPrefix(msisdn, "08") {
		return "+62" + msisdn[1:]
	} else if strings.HasPrefix(msisdn, "8") {
		return "+62" + msisdn
	}

	return "" // jika format tidak dikenali
}
