package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort        string
	DatabaseURL     string
	RedisAddr       string
	RedisDB         int
	RedisPassword   string
	ReportCacheTTL  time.Duration
	BudgetsCacheTTL time.Duration

	RedisDialTimeout  time.Duration
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
	RedisPoolTimeout  time.Duration
	RedisPoolSize     int
	RedisMinIdleConns int
}

func Load() *Config {
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	// TTL
	reportTTL := getEnvDuration("REDIS_REPORT_TTL", 30*time.Second)
	budgetsTTL := getEnvDuration("REDIS_BUDGETS_TTL", 10*time.Second)

	// Timeouts/pool
	dialTimeout := getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second)
	readTimeout := getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second)
	writeTimeout := getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second)
	poolTimeout := getEnvDuration("REDIS_POOL_TIMEOUT", 4*time.Second)
	poolSize := getEnvInt("REDIS_POOL_SIZE", 10)
	minIdleConns := getEnvInt("REDIS_MIN_IDLE_CONNS", 2)

	return &Config{
		GRPCPort:        getEnv("GRPC_PORT", "9090"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:         redisDB,
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		ReportCacheTTL:  reportTTL,
		BudgetsCacheTTL: budgetsTTL,

		RedisDialTimeout:  dialTimeout,
		RedisReadTimeout:  readTimeout,
		RedisWriteTimeout: writeTimeout,
		RedisPoolTimeout:  poolTimeout,
		RedisPoolSize:     poolSize,
		RedisMinIdleConns: minIdleConns,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
