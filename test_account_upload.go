package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// URL of the endpoint
	url := "http://localhost:8080/accounts/import" // Change port if needed

	// Create a new multipart writer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add user_id field
	userIDField, err := w.CreateFormField("user_id")
	if err != nil {
		fmt.Println("Error creating user_id field:", err)
		return
	}
	userIDField.Write([]byte("123"))

	// Create a test CSV file if it doesn't exist
	testCSVPath := "test_accounts.csv"
	if _, err := os.Stat(testCSVPath); os.IsNotExist(err) {
		f, err := os.Create(testCSVPath)
		if err != nil {
			fmt.Println("Error creating test CSV file:", err)
			return
		}
		f.WriteString("id,name,email,phone,city,account_name,account_type,account_category,account_code,contact_name,email_account,potensi,website_account,system_informasi_akademik,ownership\n")
		f.WriteString("1,John Doe,john@example.com,123456789,Jakarta,Test Account,Type A,Category B,ACC123,Contact Person,contact@example.com,High,www.example.com,SIA,Private\n")
		f.Close()
	}

	// Add the file
	f, err := os.Open(testCSVPath)
	if err != nil {
		fmt.Println("Error opening test CSV file:", err)
		return
	}
	defer f.Close()

	fileField, err := w.CreateFormFile("file_csv", filepath.Base(testCSVPath))
	if err != nil {
		fmt.Println("Error creating file field:", err)
		return
	}
	if _, err = io.Copy(fileField, f); err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// Close the writer
	w.Close()

	// Create the request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(respBody))
}
