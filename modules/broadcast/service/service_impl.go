package service

import (
	"byu-crm-service/modules/broadcast/repository"
)

type broadcastService struct {
	repo repository.BroadcastRepository
}

func NewBroadcastService(repo repository.BroadcastRepository) BroadcastService {
	return &broadcastService{repo: repo}
}
