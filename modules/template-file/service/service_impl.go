package service

import (
	"byu-crm-service/modules/template-file/repository"
	"byu-crm-service/modules/template-file/response"
)

type templateFileService struct {
	repo repository.TemplateFileRepository
}

func NewTemplateFileService(repo repository.TemplateFileRepository) TemplateFileService {
	return &templateFileService{repo: repo}
}

func (s *templateFileService) GetAllTemplateFiles(type_file string) []response.TemplateFile {
	return s.repo.GetAllTemplateFiles(type_file)
}
