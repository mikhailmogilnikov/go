package config

import (
	"os"
	"time"
)

// Config конфигурация Ledger сервиса
type Config struct {
	GRPCPort    string
	DatabaseURL string
	RedisAddr   string
	CacheTTL    time.Duration
}

// Load загружает конфигурацию
func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "9090"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		CacheTTL:    5 * time.Minute,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}



