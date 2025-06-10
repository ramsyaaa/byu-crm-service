package main

import (
	"log"

	"byu-crm-service/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables instead")
	}

	config.InitRedis()

	// Connect to database
	db := config.Connect()

	// Load and validate log configuration
	logConfig := config.LoadLogConfig()
	if err := logConfig.Validate(); err != nil {
		log.Printf("Log configuration validation failed: %v", err)
	}

	// Initialize log management with configuration
	config.InitializeLogManagement(db, logConfig)

	// Start the application
	config.Route(db)
}
