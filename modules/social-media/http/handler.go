package http

import (
	"byu-crm-service/modules/social-media/service"
)

type SocialMediaHandler struct {
	service service.SocialMediaService
}

func NewSocialMediaHandler(service service.SocialMediaService) *SocialMediaHandler {
	return &SocialMediaHandler{service: service}
}
