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

type ReportCache struct {
	Categories    []domain.CategorySummary
	TotalExpenses float64
}

type Cache struct {
	client         *redis.Client
	reportTTL      time.Duration
	budgetsListTTL time.Duration
}

type RedisClientOptions struct {
	Addr         string
	DB           int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	PoolSize     int
	MinIdleConns int
}

type TTLConfig struct {
	ReportTTL      time.Duration
	BudgetsListTTL time.Duration
}

func NewCache(redisOpts RedisClientOptions, ttlCfg TTLConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisOpts.Addr,
		DB:           redisOpts.DB,
		Password:     redisOpts.Password,
		DialTimeout:  redisOpts.DialTimeout,
		ReadTimeout:  redisOpts.ReadTimeout,
		WriteTimeout: redisOpts.WriteTimeout,
		PoolTimeout:  redisOpts.PoolTimeout,
		PoolSize:     redisOpts.PoolSize,
		MinIdleConns: redisOpts.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Printf("[REDIS] Connected to %s, DB=%d poolSize=%d minIdle=%d",
		redisOpts.Addr, redisOpts.DB, redisOpts.PoolSize, redisOpts.MinIdleConns,
	)

	return &Cache{
		client:         client,
		reportTTL:      ttlCfg.ReportTTL,
		budgetsListTTL: ttlCfg.BudgetsListTTL,
	}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

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

func (c *Cache) SetReport(ctx context.Context, userID int64, from, to time.Time, report *ReportCache) error {
	key := c.reportKey(userID, from, to)
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	log.Printf("[CACHE SET] report key=%s ttl=%v", key, c.reportTTL)
	return c.client.Set(ctx, key, data, c.reportTTL).Err()
}

func (c *Cache) InvalidateReports(ctx context.Context, userID int64) {
	pattern := c.reportSummaryPattern(userID)
	iter := c.client.Scan(ctx, 0, pattern, 200).Iterator()

	const batchSize = 200
	keys := make([]string, 0, batchSize)

	flush := func() {
		if len(keys) == 0 {
			return
		}
		if err := c.client.Unlink(ctx, keys...).Err(); err != nil {
			log.Printf("[CACHE UNLINK ERROR] report keys=%d err=%v", len(keys), err)
		} else {
			log.Printf("[CACHE UNLINK] report keys=%d", len(keys))
		}
		keys = keys[:0]
	}

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		if len(keys) >= batchSize {
			flush()
		}
	}
	flush()
	if err := iter.Err(); err != nil {
		log.Printf("[CACHE SCAN ERROR] report pattern=%s err=%v", pattern, err)
	}
}

func (c *Cache) reportKey(userID int64, from, to time.Time) string {
	return fmt.Sprintf("%s:%d:%s:%s", c.reportSummaryPrefix(), userID, from.Format("2006-01-02"), to.Format("2006-01-02"))
}

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

func (c *Cache) SetBudgets(ctx context.Context, userID int64, budgets []domain.Budget) error {
	key := c.budgetsKey(userID)
	data, err := json.Marshal(budgets)
	if err != nil {
		return err
	}
	log.Printf("[CACHE SET] budgets key=%s ttl=%v", key, c.budgetsListTTL)
	return c.client.Set(ctx, key, data, c.budgetsListTTL).Err()
}

func (c *Cache) InvalidateBudgets(ctx context.Context, userID int64) {
	key := c.budgetsKey(userID)
	if err := c.client.Unlink(ctx, key).Err(); err != nil {
		log.Printf("[CACHE UNLINK ERROR] budgets key=%s err=%v", key, err)
		return
	}
	log.Printf("[CACHE UNLINK] budgets key=%s", key)
}

func (c *Cache) budgetsKey(userID int64) string {
	return fmt.Sprintf("%s:%d", c.budgetsPrefix(), userID)
}

func (c *Cache) reportSummaryPrefix() string {
	return "report:summary"
}

func (c *Cache) reportSummaryPattern(userID int64) string {
	return fmt.Sprintf("%s:%d:*", c.reportSummaryPrefix(), userID)
}

func (c *Cache) budgetsPrefix() string {
	return "budgets:all"
}
