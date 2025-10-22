# ДЗ 5 - REST API для Ledger

REST API для управления транзакциями и бюджетами.

## Запуск

```bash
cd hw5/gateway
go run main.go
```

Сервер запустится на `http://localhost:8080`

## API Endpoints

### Транзакции

**Создать транзакцию**
```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":450,"category":"еда","description":"ланч","date":"2025-10-20"}'
```

**Получить список транзакций**
```bash
curl http://localhost:8080/api/transactions
```

### Бюджеты

**Создать бюджет**
```bash
curl -X POST http://localhost:8080/api/budgets \
-H "Content-Type: application/json" \
-d '{"category":"еда","limit":5000}'
```

**Получить список бюджетов**
```bash
curl http://localhost:8080/api/budgets
```

## Коды ответов

- `200 OK` - успешное получение данных
- `201 Created` - успешное создание
- `400 Bad Request` - ошибка валидации
- `409 Conflict` - превышен лимит бюджета
- `500 Internal Server Error` - внутренняя ошибка сервера

