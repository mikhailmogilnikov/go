# ДЗ 7 - Базы данных

REST API Ledger/Gateway с хранением данных в PostgreSQL и кэшированием отчётов в Redis.

## Подготовка БД и миграции

1) Установите Goose.
2) Задайте переменную `DATABASE_URL`, например:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
```
3) Примените миграции:
```bash
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" up
```
Откат:
```bash
goose -dir ./hw7/ledger/migrations postgres "$DATABASE_URL" down
```

Структура:
- budgets(id, category UNIQUE, limit_amount > 0)
- expenses(id, amount <> 0, category, description, date)

## Redis

Запуск:
```bash
docker run -p 6379:6379 --name cashcraft-redis -d redis:7-alpine
```
Переменные:
- `REDIS_ADDR` (по умолчанию `localhost:6379`)
- `REDIS_DB` (по умолчанию `0`)
- `REDIS_PASSWORD` (опционально)

## Запуск сервисов

Ledger инициализируется автоматически при импорте из Gateway, используя:
- `DATABASE_URL` или `DB_HOST/DB_PORT/DB_USER/DB_PASS/DB_NAME`
- Параметры пула: `SetMaxOpenConns(10)`, `SetMaxIdleConns(5)`

Запуск Gateway:
```bash
cd hw7/gateway
go run main.go
```
Сервер: `http://localhost:8080`

## Проверка через HTTP

Создать бюджет:
```bash
curl -X POST http://localhost:8080/api/budgets \
-H "Content-Type: application/json" \
-d '{"category":"еда","limit":5000}'
```
Список бюджетов:
```bash
curl http://localhost:8080/api/budgets
```
Создать транзакцию:
```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":450,"category":"еда","description":"ланч","date":"2025-10-20"}'
```
Список транзакций:
```bash
curl http://localhost:8080/api/transactions
```
Превышение бюджета:
```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":999999,"category":"еда","description":"тест","date":"2025-10-20"}'
```
Ожидается `409` и отсутствие вставки.

## Кэш отчётов

В пакете `ledger` добавлена функция `GetReportSummary(ctx, from, to)`:
- ключ `report:summary:<from>:<to>`
- TTL 30 секунд
- формат JSON

Проверка кеша через HTTP (вызовите дважды подряд, чтобы увидеть эффект кеша):
```bash
curl "http://localhost:8080/api/reports/summary?from=2025-10-01&to=2025-10-31"
```

Пример вызова через код:
```go
items, _ := ledger.GetReportSummary(context.Background(),
    time.Date(2025,10,1,0,0,0,0,time.UTC),
    time.Date(2025,10,31,0,0,0,0,time.UTC))
_ = items
```

## Примечание

Gateway остаётся без изменений по контрактам; Ledger теперь хранит данные в PostgreSQL, а отчёты кэшируются в Redis.


