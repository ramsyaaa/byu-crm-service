package service

import (
	"byu-crm-service/modules/menu/repository"
)

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}
