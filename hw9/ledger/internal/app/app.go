package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/cache"
	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/repository/pg"
	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/service"
)

func buildDSNFromEnv() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "postgres"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "cashapp"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
}

func InitService(ctx context.Context) (service.Service, func() error, error) {
	dsn := buildDSNFromEnv()
	if dsn == "" {
		return nil, nil, fmt.Errorf("empty DSN: set DATABASE_URL or DB_* variables")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Printf("ledger: connected to database")

	var cacheClient cache.CacheClient
	redisCache, err := cache.NewRedisCache(ctx)
	if err != nil {
		log.Printf("ledger: redis unavailable, continuing without cache: %v", err)
		cacheClient = &cache.NoOpCache{}
	} else {
		cacheClient = redisCache
	}

	budgetRepo := pg.NewBudgetRepository(db)
	transactionRepo := pg.NewTransactionRepository(db)

	ledgerService := service.NewLedgerService(budgetRepo, transactionRepo)

	closeFn := func() error {
		var errs []error
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("database close error: %w", err))
		}
		if err := cacheClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("cache close error: %w", err))
		}
		if len(errs) > 0 {
			return fmt.Errorf("errors during cleanup: %v", errs)
		}
		return nil
	}

	return ledgerService, closeFn, nil
}

