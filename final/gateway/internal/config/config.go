package config

import (
	"os"
	"time"
)

// Config конфигурация Gateway
type Config struct {
	HTTPPort        string
	AuthGRPCAddr    string
	LedgerGRPCAddr  string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// Load загружает конфигурацию
func Load() *Config {
	return &Config{
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		AuthGRPCAddr:    getEnv("AUTH_GRPC_ADDR", "localhost:9091"),
		LedgerGRPCAddr:  getEnv("LEDGER_GRPC_ADDR", "localhost:9090"),
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}



