package config

import (
	"os"
	"time"
)

type Config struct {
	GRPCPort    string
	DatabaseURL string
	JWTSecret   string
	TokenTTL    time.Duration
}

func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "9091"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		TokenTTL:    24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}



