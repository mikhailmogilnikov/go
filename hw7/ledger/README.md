# Ledger

Основной сервис, который работает с базой данных. Хранит всё в PostgreSQL, а отчёты кеширует в Redis.

## База данных

Сначала нужно запустить PostgreSQL. Я использовал Docker:

```bash
docker run -d --name cashapp-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16
```

Потом создаём базу:
```bash
psql "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" -c "CREATE DATABASE cashapp;"
```

## Миграции

Устанавливаем goose:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Применяем миграции:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5433/cashapp?sslmode=disable"
goose -dir ./ledger/migrations postgres "$DATABASE_URL" up
```

Откатить можно так:
```bash
goose -dir ./ledger/migrations postgres "$DATABASE_URL" down
```

Миграции создают две таблицы: `budgets` и `expenses`. Ещё есть индекс по категории и дате.

## Запуск

Ledger сам подключается к БД при импорте. Нужно только задать переменные окружения.

Для PostgreSQL:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/cashapp?sslmode=disable"
```

Или можно по отдельности:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=cashapp
```

Для Redis:
```bash
export REDIS_ADDR=localhost:6379
export REDIS_DB=0
```

Пароль для Redis не обязателен, можно не ставить.

Запускаем Redis:
```bash
docker run -d -p 6379:6379 --name cashcraft-redis redis:7-alpine
```