package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // kosongkan jika tidak pakai password
		DB:       0,                           // gunakan DB default
	})

	// Test koneksi
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Gagal terhubung ke Redis: %v", err))
	}

	fmt.Println("Redis berhasil terhubung!")
	return RedisClient
}

// Contoh fungsi untuk menyimpan dan membaca data
func SetRedis(key string, value string, expiration time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

func GetRedis(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}
