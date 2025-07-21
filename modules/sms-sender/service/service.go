package service

type SmsSenderService interface {
	CreateSms(requestBody map[string]string, rolesName []string, userRole string, territoryID int, userID int) error
	AssignSmsToUsers(requestBody map[string]string, userIDs []int) error
}
