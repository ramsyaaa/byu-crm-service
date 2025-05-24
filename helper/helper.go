package helper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func ValidateStruct(validate *validator.Validate, req interface{}, validationMessages map[string]string) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	// Ambil refleksi tipe yang tepat (harus struct, bukan pointer)
	var ref reflect.Type
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		ref = t.Elem()
	} else {
		ref = t
	}

	for _, e := range err.(validator.ValidationErrors) {
		field, ok := ref.FieldByName(e.StructField())
		var jsonKey string
		if ok {
			jsonTag := field.Tag.Get("json")
			jsonKey = strings.Split(jsonTag, ",")[0]
		} else {
			jsonKey = e.Field()
		}

		key := jsonKey + "." + e.Tag()
		msg, found := validationMessages[key]
		if !found {
			msg = e.Error()
		}
		errors[jsonKey] = msg
	}

	return errors
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

const EarthRadius = 6371000 // meters

// Function to calculate the distance between two points on the Earth using the Haversine formula
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

// Convert degrees to radians
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// function to check if a point is within a certain radius from another point
func IsWithinRadius(radius, targetLat, targetLon, mainLat, mainLon float64) bool {
	distance := haversine(mainLat, mainLon, targetLat, targetLon)
	return distance <= radius // in meters
}

func getUintPointer(m map[string]interface{}, key string) *uint {
	if val, ok := m[key].(float64); ok && val != 0 {
		temp := uint(val)
		return &temp
	}
	return nil
}

func getStringPointer(m map[string]interface{}, key string) *string {
	if val, ok := m[key].(string); ok && val != "" {
		return &val
	}
	return nil
}

func DecodeBase64Image(data string) ([]byte, string, error) {
	re := regexp.MustCompile(`^data:(image\/[a-zA-Z]+);base64,`)
	match := re.FindStringSubmatch(data)
	if len(match) != 2 {
		return nil, "", errors.New("format base64 tidak sesuai")
	}

	mimeType := match[1]
	base64Data := strings.Replace(data, match[0], "", 1)

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, "", err
	}
	return decoded, mimeType, nil
}

func DecodeBase64File(data string) ([]byte, string, error) {
	if !strings.Contains(data, ";base64,") {
		return nil, "", errors.New("format base64 tidak valid")
	}

	parts := strings.SplitN(data, ";base64,", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("data base64 tidak lengkap")
	}

	mimeType := strings.TrimPrefix(parts[0], "data:")
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", err
	}

	return decoded, mimeType, nil
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Gagal menghapus file: %s, error: %v", filePath, err)
		return err
	}
	return nil
}

func DeleteFileIfExists(path string) error {
	baseURL := os.Getenv("BASE_URL")

	// Pastikan baseURL tidak memiliki trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	// Hapus baseURL dari awal path jika ada
	path = strings.TrimPrefix(path, baseURL)
	path = strings.TrimPrefix(path, "/") // Hapus slash di awal jika ada

	// Ubah semua backslash menjadi slash
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	// Cek apakah file ada
	if _, err := os.Stat(normalizedPath); err == nil {
		// File ditemukan, hapus
		err := os.Remove(normalizedPath)
		if err != nil {
			return fmt.Errorf("gagal menghapus file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		// Error selain file tidak ditemukan
		return fmt.Errorf("gagal mengecek file: %v", err)
	}

	// Jika file tidak ada, tidak dianggap error
	return nil
}

func SaveValidatedBase64File(base64Str string, uploadDir string) (string, error) {
	// Decode base64
	parts := strings.SplitN(base64Str, ",", 2)
	if len(parts) != 2 {
		return "", errors.New("format base64 tidak valid")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("gagal mendecode base64: " + err.Error())
	}

	// Validasi ukuran maksimum 5MB
	if len(data) > 5*1024*1024 {
		return "", errors.New("ukuran file maksimum adalah 5MB")
	}

	// Deteksi MIME
	mimeType := httpDetectContentType(data)

	// Map ekstensi berdasarkan MIME
	mimeExtensions := map[string]string{
		"image/jpeg":                    ".jpg",
		"image/png":                     ".png",
		"application/pdf":               ".pdf",
		"application/zip":               ".zip",
		"application/msword":            ".doc",
		"application/vnd.ms-excel":      ".xls",
		"application/vnd.ms-powerpoint": ".ppt",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
		"text/plain":       ".txt",
		"application/json": ".json",
	}

	// Tentukan ekstensi
	ext, ok := mimeExtensions[mimeType]
	if !ok {
		ext = ".bin"
	}

	// Buat direktori jika belum ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("gagal membuat direktori: " + err.Error())
	}

	// Simpan file
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", errors.New("gagal menyimpan file: " + err.Error())
	}

	return filePath, nil
}

func SaveUploadedFile(c *fiber.Ctx, file *multipart.FileHeader, folder string) (string, error) {
	// Buat folder tujuan jika belum ada
	dirPath := filepath.Join("public", "uploads", folder)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
		return "", fmt.Errorf("gagal membuat folder: %v", err)
	}

	// Format nama file unik
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// If no extension, try to determine from content type
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening file to detect type: %v", err)
			// Default to .bin if we can't determine
			ext = ".bin"
		} else {
			defer src.Close()

			// Read a bit of the file to detect content type
			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			if err != nil && err != io.EOF {
				log.Printf("Error reading file to detect type: %v", err)
				ext = ".bin"
			} else {
				mimeType := http.DetectContentType(buffer)

				// Map common mime types to extensions
				mimeExtensions := map[string]string{
					"image/jpeg":      ".jpg",
					"image/png":       ".png",
					"image/gif":       ".gif",
					"application/pdf": ".pdf",
					"text/plain":      ".txt",
				}

				if extFromMime, ok := mimeExtensions[mimeType]; ok {
					ext = extFromMime
				} else {
					ext = ".bin"
				}

				// Reset file pointer for saving
				src.Seek(0, 0)
			}
		}
	}

	fileName := fmt.Sprintf("%d%s", timestamp, ext)
	fullPath := filepath.Join(dirPath, fileName)

	// Try to save the file with better error handling
	if err := c.SaveFile(file, fullPath); err != nil {
		log.Printf("Error saving file: %v", err)

		// Try alternative method if SaveFile fails
		src, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("gagal membuka file: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(fullPath)
		if err != nil {
			return "", fmt.Errorf("gagal membuat file: %v", err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return "", fmt.Errorf("gagal menyalin file: %v", err)
		}
	}

	return fullPath, nil
}

// Fungsi bantu deteksi MIME dari data file
func httpDetectContentType(data []byte) string {
	return http.DetectContentType(data[:512])
}

func CapitalizeWords(input string) string {
	input = strings.TrimSpace(input)
	words := strings.Fields(input)

	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToUpper(runes[j])
			}
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}

func UppercaseTrim(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}
