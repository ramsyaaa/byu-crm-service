package service

import (
	"byu-crm-service/modules/template-file/response"
)

type TemplateFileService interface {
	GetAllTemplateFiles(type_file string) []response.TemplateFile
}
