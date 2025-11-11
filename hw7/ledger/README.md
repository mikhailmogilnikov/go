# Ledger (hw7)

Хранилище данных для приложения: работа с PostgreSQL и кэшированием отчётов в Redis. Пакет инициализируется при импорте из Gateway: устанавливает соединения с PostgreSQL и Redis, применяет бизнес‑логику (валидация, лимиты бюджетов) и предоставляет функции `SetBudget`, `ListBudgets`, `AddTransaction`, `ListTransactions`, `GetReportSummary`.

## Переменные окружения

- `DATABASE_URL` — строка подключения к PostgreSQL. Пример: `postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable`.
- Если `DATABASE_URL` не задана, используется сборка из переменных: `DB_HOST` (по умолчанию `localhost`), `DB_PORT` (`5432`), `DB_USER` (`postgres`), `DB_PASS` (`postgres`), `DB_NAME` (`cashapp`), `DB_SSLMODE` (`disable`).
- Redis:
  - `REDIS_ADDR` (по умолчанию `localhost:6379`)
  - `REDIS_DB` (по умолчанию `0`)
  - `REDIS_PASSWORD` (опционально)

## Создание БД и применение миграций (Goose)

1) Установите Goose (см. документацию Goose для вашей ОС).

2) Экспортируйте `DATABASE_URL` или задайте переменные `DB_*` (см. выше).

3) Примените миграции (и откатите при необходимости):

```bash
# применить
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" up

# откатить
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" down
```

Структура схемы:

- `budgets(id, category UNIQUE, limit_amount > 0)`
- `expenses(id, amount <> 0, category, description, date)`

## Что делает Ledger

- Подключается к PostgreSQL (`sql.Open("pgx", dsn)` + `Ping`) и логирует успешное подключение.
- Настраивает пул соединений: `SetMaxOpenConns(10)`, `SetMaxIdleConns(5)`, `SetConnMaxLifetime(...)`.
- Реализует CRUD через БД, в том числе:
  - `SetBudget` — upsert бюджета категории (`ON CONFLICT (category) DO UPDATE`), инвалидирует кеш `budgets:all`.
  - `ListBudgets` — чтение из таблицы `budgets` (опционально возвращает результат из Redis‑кеша с коротким TTL).
  - `AddTransaction` — проверка лимита по категории и вставка в `expenses`; при превышении возвращает ошибку `budget exceeded`.
  - `ListTransactions` — чтение из `expenses`.
- Инициализирует Redis‑клиент и логирует успешное подключение.
- Кэширует результат `GetReportSummary(ctx, from, to)` в Redis:
  - ключ `report:summary:<from>:<to>` (даты в формате `YYYY-MM-DD`);
  - формат значения — JSON;
  - TTL по умолчанию — 30 секунд.

## Примечание о запуске

Ledger не поднимает собственный HTTP‑сервер в этом задании. Он используется Gateway через импорт пакета. Перед запуском Gateway убедитесь, что миграции применены и переменные окружения для БД/Redis заданы (см. README Gateway).


