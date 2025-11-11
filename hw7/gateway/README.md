# Gateway (hw7)

REST‑обёртка над Ledger. Эндпоинты те же, что были раньше, плюс простой отчёт.

## Запуск

1) Примените миграции (см. `hw7/ledger/README.md`).  
2) Задайте переменные (Ledger читает их при инициализации):

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export REDIS_DB="0"
```

3) Запустите сервер:
```bash
cd hw7/gateway
go run main.go
# http://localhost:8080
```

## Быстрые проверки (cURL)

Создать бюджет:
```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"еда","limit":5000}'
```

Создать транзакцию:
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":450,"category":"еда","description":"ланч","date":"2025-10-20"}'
```

Списки:
```bash
curl http://localhost:8080/api/budgets
curl http://localhost:8080/api/transactions
```

Превышение (ожидается 409):
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":999999,"category":"еда","description":"тест","date":"2025-10-20"}'
```

Отчёт (вызовите два раза подряд, второй быстрее из‑за кеша):
```bash
curl "http://localhost:8080/api/reports/summary?from=2025-10-01&to=2025-10-31"
curl "http://localhost:8080/api/reports/summary?from=2025-10-01&to=2025-10-31"
```

Теперь Ledger хранит данные в PostgreSQL, а отчёты кэшируются в Redis на короткий TTL.


