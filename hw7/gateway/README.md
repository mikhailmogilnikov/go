# Gateway

HTTP сервер для тестирования API. Просто обращается к Ledger, который уже работает с PostgreSQL и Redis.

## Запуск

Сначала запустим PostgreSQL и Redis (смотрим в ledger/README.md как это сделать).

Потом нужно задать переменные:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5433/cashapp?sslmode=disable"
export REDIS_ADDR=localhost:6379
export REDIS_DB=0
```

И запускаем:
```bash
cd gateway
go run .
```

Сервер будет на `http://localhost:8080`

## Примеры запросов

### Бюджеты

Создать:
```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"food","limit":1000}'
```

Получить все:
```bash
curl http://localhost:8080/api/budgets
```

### Транзакции

Создать:
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":200,"category":"food","description":"groceries","date":"2025-11-11"}'
```

Получить все:
```bash
curl http://localhost:8080/api/transactions
```

Попробовать превысить лимит (должно вернуть 409):
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":1500,"category":"food","description":"big expense","date":"2025-11-11"}'
```

### Отчёты

Первый запрос (пойдёт в БД):
```bash
curl "http://localhost:8080/api/reports/summary?from=2025-11-10&to=2025-11-11"
```

Второй запрос (из кеша, будет быстрее):
```bash
curl "http://localhost:8080/api/reports/summary?from=2025-11-10&to=2025-11-11"
```

Второй раз быстрее, потому что данные уже в Redis. Кеш живёт 30 секунд.