# ДЗ 8: Чистая архитектура

Рефакторинг проекта под чистую архитектуру - разделил код на слои.

## Структура

### Ledger

```
ledger/
├── internal/
│   ├── domain/          # Сущности и интерфейсы (Transaction, Budget)
│   ├── repository/pg/   # Работа с PostgreSQL
│   ├── service/         # Бизнес-логика
│   ├── app/             # Инициализация (InitService)
│   └── cache/           # Redis
├── migrations/
└── go.mod
```

### Gateway

```
gateway/
├── internal/api/   # HTTP handlers
├── main.go
└── go.mod
```

## Как работает

1. **Domain** - сущности и интерфейсы, ни от чего не зависит
2. **Repository** - реализует интерфейсы из domain, работает с БД
3. **Service** - бизнес-логика (валидация, проверка бюджета)
4. **App** - собирает все вместе через `InitService()`
5. **Gateway** - HTTP слой, знает только про интерфейс сервиса

Зависимости идут внутрь: Gateway → Service → Domain ← Repository

## Запуск

1. Настроить переменные окружения:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
```

2. Выполнить миграции:
```bash
cd ledger
goose -dir migrations postgres "$DATABASE_URL" up
```

3. Запустить сервер:
```bash
cd gateway
go run main.go
```

## API

- `POST /api/budgets` - создать бюджет
- `GET /api/budgets` - получить бюджеты
- `POST /api/transactions` - создать транзакцию
- `GET /api/transactions` - получить транзакции
- `GET /api/reports/summary?from=YYYY-MM-DD&to=YYYY-MM-DD` - отчет

## Что изменилось

Раньше все было в одном месте. Теперь:
- Domain не зависит ни от чего
- Gateway не знает про БД
- Бизнес-логика в сервисе
- Все собирается через фабрику `InitService()`
