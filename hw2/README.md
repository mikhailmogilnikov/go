# Домашнее задание 2 - Функции в Go

Проект содержит два сервиса: **Gateway** (HTTP-шлюз) и **Ledger** (бизнес-логика для работы с транзакциями).

## Задание 1: Gateway и Ledger

**Запуск:**

```bash
cd gateway
go run main.go
```

**Проверка:**

В другом терминале выполните:

```bash
curl http://localhost:8080/ping
```

Ожидаемый ответ: `pong`

Или откройте в браузере: http://localhost:8080/ping

---

### Ledger сервис

**Запуск:**

```bash
cd ledger
go run main.go
```

## Задание 2: Структура Transaction и функции

## Запуск обоих сервисов одновременно

**Терминал 1 (Gateway):**
```bash
cd gateway
go run main.go
```

**Терминал 2 (Ledger):**
```bash
cd ledger
go run main.go
```

**Терминал 3 (Проверка Gateway):**
```bash
curl http://localhost:8080/ping
```

## Требования

- Go 1.21 или выше
- Доступ к портам 8080 (для Gateway)
