package repository

import (
	"byu-crm-service/modules/template-file/response"
	"os"

	"gorm.io/gorm"
)

type templateFileRepository struct {
	db *gorm.DB
}

func NewTemplateFileRepository(db *gorm.DB) TemplateFileRepository {
	return &templateFileRepository{db: db}
}

var template_file_lists = []response.TemplateFile{
	{Name: "Tutorial YAE", Type: "tutorial-yae", FilePath: "public/format/Role-YAE-Guide.pdf"},
	{Name: "Tutorial Branch", Type: "tutorial-branch", FilePath: "public/format/Role-Branch-Guide.pdf"},
	{Name: "Tutorial Regional", Type: "tutorial-regional", FilePath: "public/format/Role-Regional-Guide.pdf"},
	{Name: "Tutorial Area", Type: "tutorial-area", FilePath: "public/format/Role-Area-Guide.pdf"},
}

// Function untuk filter sesuai type
func (r *templateFileRepository) GetAllTemplateFiles(type_file string) []response.TemplateFile {
	baseURL := os.Getenv("BASE_URL")

	// hasil akhir
	var result []response.TemplateFile

	if type_file == "" {
		for _, f := range template_file_lists {
			f.FilePath = baseURL + f.FilePath
			result = append(result, f)
		}
		return result
	}

	// filter berdasarkan type
	for _, f := range template_file_lists {
		if f.Type == type_file {
			f.FilePath = baseURL + f.FilePath
			result = append(result, f)
		}
	}
	return result
}
