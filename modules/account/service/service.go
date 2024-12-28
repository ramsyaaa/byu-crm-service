package service

type AccountService interface {
	ProcessAccount(data []string) error
}
