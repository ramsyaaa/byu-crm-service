package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func ErrorValidationFormat(err error, validationMessages map[string]string) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		// Buat key yang sesuai dengan field dan tag error
		key := e.Field() + "." + e.Tag()
		if message, exists := validationMessages[key]; exists {
			errors[e.Field()] = message
		} else {
			errors[e.Field()] = e.Error()
		}
	}

	return errors
}

func saveFileToLocal(file *multipart.FileHeader, directory string, allowedFormats []string) (*string, error) {
	// Validate file type
	ext := filepath.Ext(file.Filename)
	ext = strings.ToLower(ext)

	// Check if the file extension is allowed
	isValidExt := false
	for _, allowedExt := range allowedFormats {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return nil, fmt.Errorf("invalid file format. Allowed formats are: %v", allowedFormats)
	}

	// Define a unique file name
	filename := fmt.Sprintf("%s%s", generateUniqueID(), ext)

	// Define the full path where to save the file
	savePath := filepath.Join("public", directory, filename)

	// Create the directory if it doesn't exist
	dir := filepath.Dir(savePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directories for file storage: %v", err)
		}
	}

	// Open the file from the incoming multipart request
	fileSrc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer fileSrc.Close()

	// Create the destination file on the server
	fileDest, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer fileDest.Close()

	// Copy the content of the uploaded file to the destination file
	_, err = io.Copy(fileDest, fileSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Return the relative path to the saved file
	filePath := fmt.Sprintf("/%s/%s", directory, filename)
	return &filePath, nil
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
