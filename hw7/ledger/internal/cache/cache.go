package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	initOnce    sync.Once
)

func buildRedisConfig() *redis.Options {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			db = parsed
		}
	}

	password := os.Getenv("REDIS_PASSWORD")

	return &redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	}
}

func Init() {
	var initErr error

	initOnce.Do(func() {
		config := buildRedisConfig()
		client := redis.NewClient(config)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("redis ping failed: %w", err)
			return
		}

		redisClient = client
		log.Printf("ledger: connected to redis at %s", config.Addr)
	})

	if initErr != nil {
		log.Fatalf("ledger: redis init failed: %v", initErr)
	}
}

func Client() *redis.Client {
	return redisClient
}

func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
