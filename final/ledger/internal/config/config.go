package config

import (
	"os"
	"strconv"
	"time"
)

// Config конфигурация Ledger сервиса
type Config struct {
	GRPCPort      string
	DatabaseURL   string
	RedisAddr     string
	RedisDB       int
	RedisPassword string
	CacheTTL      time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	return &Config{
		GRPCPort:      getEnv("GRPC_PORT", "9090"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:       redisDB,
		RedisPassword: os.Getenv("REDIS_PASSWORD"), // опционально
		CacheTTL:      30 * time.Second,            // TTL 30 секунд как в ТЗ
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
