# Ledger (hw7)

Хранит данные в PostgreSQL и кэширует отчёты в Redis. Пакет инициализируется при импорте из Gateway.

## Переменные окружения

- `DATABASE_URL` (пример: `postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable`)
- или отдельно: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_SSLMODE`
- Redis: `REDIS_ADDR` (по умолчанию `localhost:6379`), `REDIS_DB` (`0`), `REDIS_PASSWORD` (опционально)

## Миграции (Goose)

```bash
# применить
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" up
# откатить
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" down
```

Таблицы: `budgets(category UNIQUE, limit_amount>0)`, `expenses(amount<>0, category, description, date)`.

## Что внутри (кратко)

- Подключение к PostgreSQL (`sql.Open("pgx", dsn)`, `Ping`), пул: open=10, idle=5.
- `SetBudget` — upsert; `ListBudgets` — SELECT (+ опциональный кеш).
- `AddTransaction` — проверка лимита категории, затем INSERT.
- `ListTransactions` — SELECT.
- `GetReportSummary(from,to)` — агрегирует траты по категориям и кладёт результат в Redis на ~30 секунд.

Перед запуском Gateway примените миграции и задайте переменные из этого файла.


