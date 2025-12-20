package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
)

// ReportCache структура для кэширования отчётов
type ReportCache struct {
	Categories    []domain.CategorySummary
	TotalExpenses float64
}

// Cache обёртка над Redis
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCache создаёт новый кэш
func NewCache(addr string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Cache{
		client: client,
		ttl:    ttl,
	}, nil
}

// Close закрывает соединение
func (c *Cache) Close() error {
	return c.client.Close()
}

// GetReport получает отчёт из кэша
func (c *Cache) GetReport(ctx context.Context, userID int64, from, to time.Time) (*ReportCache, error) {
	key := c.reportKey(userID, from, to)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("[CACHE MISS] report key=%s", key)
			return nil, nil
		}
		return nil, err
	}

	var report ReportCache
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	log.Printf("[CACHE HIT] report key=%s", key)
	return &report, nil
}

// SetReport сохраняет отчёт в кэш
func (c *Cache) SetReport(ctx context.Context, userID int64, from, to time.Time, report *ReportCache) error {
	key := c.reportKey(userID, from, to)
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	log.Printf("[CACHE SET] report key=%s ttl=%v", key, c.ttl)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// InvalidateReports удаляет кэш отчётов пользователя
func (c *Cache) InvalidateReports(ctx context.Context, userID int64) {
	pattern := fmt.Sprintf("report:%d:*", userID)
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
}

func (c *Cache) reportKey(userID int64, from, to time.Time) string {
	return fmt.Sprintf("report:%d:%s:%s", userID, from.Format("2006-01-02"), to.Format("2006-01-02"))
}

