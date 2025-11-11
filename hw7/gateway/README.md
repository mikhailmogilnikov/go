# Gateway (hw7)

HTTP‑сервис для работы с Ledger через REST API. Конtrakты эндпоинтов прежние: бюджеты и транзакции. Также добавлен эндпоинт для проверки отчёта и кэша.

## Требования

- PostgreSQL доступен и миграции Ledger применены (см. `hw7/ledger/README.md`).
- Заданы переменные окружения для подключения Ledger к PostgreSQL и Redis (Gateway передаёт их процессу, Ledger читает при инициализации):
  - `DATABASE_URL` или `DB_HOST/DB_PORT/DB_USER/DB_PASS/DB_NAME/DB_SSLMODE`
  - `REDIS_ADDR`, `REDIS_DB`, `REDIS_PASSWORD` (опционально)

## Запуск

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export REDIS_DB="0"

cd hw7/gateway
go run main.go
# сервер слушает http://localhost:8080
```

## Эндпоинты и примеры cURL

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

Проверка превышения лимита (ожидается 409):
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":999999,"category":"еда","description":"тест","date":"2025-10-20"}'
```

Отчёт с кэшем (вызвать дважды подряд: сначала cache miss, затем cache hit):
```bash
curl "http://localhost:8080/api/reports/summary?from=2025-10-01&to=2025-10-31"
curl "http://localhost:8080/api/reports/summary?from=2025-10-01&to=2025-10-31"
```

## Примечание

Gateway обращается к тем же HTTP‑эндпоинтам, которые вы использовали ранее, но теперь Ledger хранит данные в PostgreSQL, а результаты отчётов кэшируются в Redis на короткий TTL.


