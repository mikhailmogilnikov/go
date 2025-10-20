# ДЗ 4 - ООП в Go

## Что сделано

1. Создан интерфейс `Validatable` с методом `Validate() error`
2. Реализована валидация для `Transaction`:
   - Amount > 0
   - Category не пустая
3. Реализована валидация для `Budget`:
   - Limit > 0
   - Category не пустая
4. Добавлена функция `CheckValid(v Validatable) error` для демонстрации полиморфизма
5. Валидация встроена в `AddTransaction()` и `SetBudget()`

## Запуск

```bash
cd hw4/ledger
go run main.go
```

Программа загружает бюджеты из файла, тестирует валидацию разных объектов через интерфейс и добавляет транзакции с проверкой.
