package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS budgets (
			id SERIAL PRIMARY KEY,
			category TEXT UNIQUE NOT NULL,
			limit_amount NUMERIC(14,2) NOT NULL CHECK (limit_amount > 0)
		)`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			amount NUMERIC(14,2) NOT NULL CHECK (amount <> 0),
			category TEXT NOT NULL,
			description TEXT,
			date DATE NOT NULL
		)`,
	}

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Миграция %d: таблица уже существует\n", i+1)
			} else {
				log.Fatalf("Ошибка миграции %d: %v", i+1, err)
			}
		} else {
			fmt.Printf("Миграция %d: применена успешно\n", i+1)
		}
	}

	fmt.Println("Все миграции применены")
}

