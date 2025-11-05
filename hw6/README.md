# ДЗ 6 - Тестирование

REST API для управления транзакциями и бюджетами с тестами.

## Запуск сервера

```bash
cd hw6/gateway
go run main.go
```

Сервер запустится на `http://localhost:8080`

## Запуск тестов

### Unit-тесты для ledger

```bash
cd hw6/ledger
go test -v
```

### Интеграционные тесты для gateway

```bash
cd hw6/gateway
go test ./internal/api -v
```

### Запуск всех тестов

```bash
cd hw6/ledger
go test -cover

cd ../gateway
go test ./... -cover
```

## Покрытие кода

Для генерации отчёта о покрытии:

```bash
cd hw6/ledger
go test -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html

cd ../gateway
go test ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
```

**Покрытие тестами ~82%**

- ledger: 79.5% покрытия
- gateway/internal/api: 85.3% покрытия

## Структура тестов

### Unit-тесты (ledger/ledger_test.go)

- Валидация транзакций (валидные, нулевые/отрицательные суммы, пустые категории)
- Валидация бюджетов (валидные, нулевые/отрицательные лимиты, пустые категории)
- Бизнес-правила (добавление транзакций в пределах бюджета и превышение лимита)

### Интеграционные тесты (gateway/internal/api/handlers_test.go)

- Создание бюджета (успешное создание, валидация, обработка ошибок)
- Создание транзакций (успешное создание, превышение бюджета, валидация)
- Получение списков (бюджеты и транзакции)
- Интеграционный поток (создание бюджета и транзакции)

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

