package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	sqlDB    *sql.DB
	initOnce sync.Once
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

func Init() {
	var initErr error

	initOnce.Do(func() {
		dsn := buildDSNFromEnv()
		if dsn == "" {
			initErr = fmt.Errorf("empty DSN: set DATABASE_URL or DB_* variables")
			return
		}

		db, err := sql.Open("pgx", dsn)
		if err != nil {
			initErr = fmt.Errorf("sql.Open failed: %w", err)
			return
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err = db.Ping(); err != nil {
			_ = db.Close()
			initErr = fmt.Errorf("database ping failed: %w", err)
			return
		}

		sqlDB = db
		log.Printf("ledger: connected to database")
	})

	if initErr != nil {
		log.Fatalf("ledger: database init failed: %v", initErr)
	}
}

func DB() *sql.DB {
	return sqlDB
}

func Close() error {
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}
