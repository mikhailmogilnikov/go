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

// Cache обёртка над Redis для кэширования данных
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCache создаёт новый кэш с подключением к Redis
func NewCache(addr string, db int, password string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})

	// Проверяем подключение к Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Printf("[REDIS] Connected to %s, DB=%d", addr, db)

	return &Cache{
		client: client,
		ttl:    ttl,
	}, nil
}

// Close закрывает соединение с Redis
func (c *Cache) Close() error {
	return c.client.Close()
}

// === КЭШИРОВАНИЕ ОТЧЁТОВ ===

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

// SetReport сохраняет отчёт в кэш с TTL
func (c *Cache) SetReport(ctx context.Context, userID int64, from, to time.Time, report *ReportCache) error {
	key := c.reportKey(userID, from, to)
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	log.Printf("[CACHE SET] report key=%s ttl=%v", key, c.ttl)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// InvalidateReports удаляет кэш отчётов пользователя при добавлении транзакции
func (c *Cache) InvalidateReports(ctx context.Context, userID int64) {
	pattern := fmt.Sprintf("report:%d:*", userID)
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
		log.Printf("[CACHE DEL] report key=%s", iter.Val())
	}
}

func (c *Cache) reportKey(userID int64, from, to time.Time) string {
	// Формат ключа: report:summary:<userID>:<from>:<to>
	return fmt.Sprintf("report:summary:%d:%s:%s", userID, from.Format("2006-01-02"), to.Format("2006-01-02"))
}

// КЭШИРОВАНИЕ БЮДЖЕТОВ

// GetBudgets получает список бюджетов из кэша
func (c *Cache) GetBudgets(ctx context.Context, userID int64) ([]domain.Budget, error) {
	key := c.budgetsKey(userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("[CACHE MISS] budgets key=%s", key)
			return nil, nil
		}
		return nil, err
	}

	var budgets []domain.Budget
	if err := json.Unmarshal(data, &budgets); err != nil {
		return nil, err
	}
	log.Printf("[CACHE HIT] budgets key=%s", key)
	return budgets, nil
}

// SetBudgets сохраняет список бюджетов в кэш
func (c *Cache) SetBudgets(ctx context.Context, userID int64, budgets []domain.Budget) error {
	key := c.budgetsKey(userID)
	data, err := json.Marshal(budgets)
	if err != nil {
		return err
	}
	// TTL для бюджетов короче - 10 секунд
	ttl := 10 * time.Second
	log.Printf("[CACHE SET] budgets key=%s ttl=%v", key, ttl)
	return c.client.Set(ctx, key, data, ttl).Err()
}

// InvalidateBudgets удаляет кэш бюджетов после изменения
func (c *Cache) InvalidateBudgets(ctx context.Context, userID int64) {
	key := c.budgetsKey(userID)
	c.client.Del(ctx, key)
	log.Printf("[CACHE DEL] budgets key=%s", key)
}

func (c *Cache) budgetsKey(userID int64) string {
	// Формат ключа: budgets:<userID>
	return fmt.Sprintf("budgets:%d", userID)
}
