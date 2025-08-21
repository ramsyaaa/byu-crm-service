package repository

import (
	"byu-crm-service/modules/template-file/response"
)

type TemplateFileRepository interface {
	GetAllTemplateFiles(type_file string) []response.TemplateFile
}
