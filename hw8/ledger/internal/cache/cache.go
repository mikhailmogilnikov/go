package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheClient представляет интерфейс для работы с кэшем
type CacheClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

// RedisCache реализует CacheClient для Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache создает новый экземпляр RedisCache
func NewRedisCache(ctx context.Context) (*RedisCache, error) {
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

	config := &redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	}

	client := redis.NewClient(config)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Printf("ledger: connected to redis at %s", config.Addr)

	return &RedisCache{client: client}, nil
}

// Get возвращает значение по ключу
func (r *RedisCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

// Set устанавливает значение с TTL
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

// Del удаляет ключи
func (r *RedisCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

// Close закрывает соединение
func (r *RedisCache) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// NoOpCache реализует CacheClient без реального кэширования (для случаев, когда Redis недоступен)
type NoOpCache struct{}

// Get всегда возвращает ошибку
func (n *NoOpCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringResult("", redis.Nil)
}

// Set ничего не делает
func (n *NoOpCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("OK", nil)
}

// Del ничего не делает
func (n *NoOpCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return redis.NewIntResult(0, nil)
}

// Close ничего не делает
func (n *NoOpCache) Close() error {
	return nil
}

