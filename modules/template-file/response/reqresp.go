package response

type TemplateFile struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}
