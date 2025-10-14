# Домашнее задание 3

Сервис Ledger с поддержкой бюджетов.

## Что реализовано

1. Структура `Budget` с полями Category, Limit, Period
2. Хранилище бюджетов `map[string]Budget`
3. Функция `AddTransaction` с проверкой бюджета
4. Функция `SetBudget` для установки бюджетов
5. Функция `LoadBudgets` для загрузки из JSON файла
6. Обработка ошибок при превышении бюджета

## Запуск

```bash
cd hw3/ledger
go run main.go
```

## Структура

```
hw3/
├── README.md
└── ledger/
    ├── main.go
    ├── go.mod
    └── budgets.json
```

## Требования

Go 1.21+
