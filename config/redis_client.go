package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var Ctx = context.Background()

// Redis Client global
var RedisClient *redis.Client

// Init Redis
func InitRedis() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	// Inisialisasi Redis Client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // Kosong jika tanpa password
		DB:       redisDB,       // Gunakan database Redis ke-0
	})

	// Cek koneksi Redis
	pong, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis Connection Failed: %v", err)
	}
	fmt.Println("✅ Redis Connected:", pong)
}
