#!/bin/bash

set -e

echo "=== Тестирование hw7 ==="

cd "$(dirname "$0")"

export DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable}"
export REDIS_ADDR="${REDIS_ADDR:-localhost:6379}"
export REDIS_DB="${REDIS_DB:-0}"

echo "DATABASE_URL: $DATABASE_URL"
echo "REDIS_ADDR: $REDIS_ADDR"
echo ""
echo "💡 Если нужно изменить DATABASE_URL, установите переменную окружения:"
echo "   export DATABASE_URL='postgres://user:password@host:port/database?sslmode=disable'"
echo ""

echo ""
echo "=== Шаг 1: Проверка зависимостей ==="

if ! command -v docker &> /dev/null; then
    echo "Ошибка: docker не установлен"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Ошибка: go не установлен"
    exit 1
fi

echo "✓ Все зависимости установлены"

echo ""
echo "=== Шаг 2: Настройка PostgreSQL ==="
PG_PORT=5433
if lsof -ti:$PG_PORT > /dev/null 2>&1; then
    echo "Порт $PG_PORT занят, проверяю контейнер..."
    PG_PORT=5434
fi

if docker ps | grep -q cashcraft-pg; then
    echo "✓ PostgreSQL контейнер уже запущен"
    CONTAINER_PORT=$(docker port cashcraft-pg 5432/tcp 2>/dev/null | cut -d: -f2)
    if [ -n "$CONTAINER_PORT" ]; then
        PG_PORT=$CONTAINER_PORT
    fi
elif docker ps -a | grep -q cashcraft-pg; then
    echo "Запуск существующего PostgreSQL контейнера..."
    docker start cashcraft-pg
    echo "Ожидание запуска PostgreSQL..."
    sleep 5
    CONTAINER_PORT=$(docker port cashcraft-pg 5432/tcp 2>/dev/null | cut -d: -f2)
    if [ -n "$CONTAINER_PORT" ]; then
        PG_PORT=$CONTAINER_PORT
    fi
else
    echo "Создание нового PostgreSQL контейнера на порту $PG_PORT..."
    echo "  Пользователь: postgres"
    echo "  Пароль: postgres"
    echo "  База данных: cashapp"
    echo "  Порт: $PG_PORT"
    docker rm -f cashcraft-pg 2>/dev/null
    docker run --name cashcraft-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=cashapp -p $PG_PORT:5432 -d postgres:16
    echo "Ожидание запуска PostgreSQL..."
    sleep 5
fi

if ! docker ps | grep -q cashcraft-pg; then
    echo "⚠ Ошибка: PostgreSQL контейнер не запущен"
    exit 1
fi

echo "✓ PostgreSQL готов на порту $PG_PORT"
export DATABASE_URL="postgres://postgres:postgres@localhost:$PG_PORT/cashapp?sslmode=disable"
echo "DATABASE_URL установлен: $DATABASE_URL"

echo ""
echo "=== Шаг 3: Запуск Redis (если не запущен) ==="
if docker ps | grep -q cashcraft-redis; then
    echo "Redis контейнер уже запущен"
elif docker ps -a | grep -q cashcraft-redis; then
    echo "Запуск существующего Redis контейнера..."
    docker start cashcraft-redis
    echo "Ожидание запуска Redis..."
    sleep 2
else
    echo "Создание нового Redis контейнера..."
    docker run -p 6379:6379 --name cashcraft-redis -d redis:7-alpine
    echo "Ожидание запуска Redis..."
    sleep 2
fi

echo ""
echo "=== Шаг 4: Применение миграций ==="
echo "Применение миграций через Go скрипт..."
cd scripts
go mod tidy
go run apply_migrations.go
cd ..

echo ""
echo "Очистка старых данных..."
cd scripts
cat > clear_db.go <<'CLEAREOF'
package main
import (
	"database/sql"
	"log"
	"os"
	_ "github.com/jackc/pgx/v5/stdlib"
)
func main() {
	db, _ := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	defer db.Close()
	db.Exec("TRUNCATE TABLE expenses, budgets RESTART IDENTITY CASCADE")
	log.Println("База данных очищена")
}
CLEAREOF
go run clear_db.go
rm clear_db.go
cd ..

echo ""
echo "=== Шаг 5: Сборка gateway ==="
cd gateway
go build -o gateway .
cd ..

echo ""
echo "=== Шаг 6: Запуск gateway в фоне ==="
pkill -f "gateway" 2>/dev/null || true
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 2
cd gateway
export DATABASE_URL="$DATABASE_URL"
export REDIS_ADDR="$REDIS_ADDR"
export REDIS_DB="$REDIS_DB"
./gateway > /tmp/gateway.log 2>&1 &
GATEWAY_PID=$!
cd ..
echo "Gateway запущен (PID: $GATEWAY_PID) с DATABASE_URL=$DATABASE_URL"

echo "Ожидание запуска сервера..."
sleep 3
if ! kill -0 $GATEWAY_PID 2>/dev/null; then
    echo "ОШИБКА: Gateway не запустился"
    tail -10 /tmp/gateway.log
    exit 1
fi

cleanup() {
    echo ""
    echo "=== Остановка gateway ==="
    kill $GATEWAY_PID 2>/dev/null || true
    wait $GATEWAY_PID 2>/dev/null || true
}

trap cleanup EXIT

echo ""
echo "=== Шаг 7: Тестирование API ==="

BASE_URL="http://localhost:8080"

echo ""
echo "7.1. Создание бюджета..."
BUDGET_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/api/budgets" \
  -H "Content-Type: application/json" \
  -d '{"category":"еда","limit":5000}')
HTTP_CODE=$(echo "$BUDGET_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BUDGET_BODY=$(echo "$BUDGET_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $BUDGET_BODY"
if [ "$HTTP_CODE" != "201" ]; then
    echo "ОШИБКА: Ожидался код 201"
    exit 1
fi

echo ""
echo "7.2. Получение списка бюджетов..."
BUDGETS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/api/budgets")
HTTP_CODE=$(echo "$BUDGETS_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BUDGETS_BODY=$(echo "$BUDGETS_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $BUDGETS_BODY"
if [ "$HTTP_CODE" != "200" ]; then
    echo "ОШИБКА: Ожидался код 200"
    exit 1
fi

echo ""
echo "7.3. Создание транзакции (в пределах лимита)..."
TX1_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/api/transactions" \
  -H "Content-Type: application/json" \
  -d '{"amount":450,"category":"еда","description":"ланч","date":"2025-10-20"}')
HTTP_CODE=$(echo "$TX1_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
TX1_BODY=$(echo "$TX1_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $TX1_BODY"
if [ "$HTTP_CODE" != "201" ]; then
    echo "ОШИБКА: Ожидался код 201"
    exit 1
fi

echo ""
echo "7.4. Создание ещё одной транзакции (450 + 3000 = 3450 < 5000)..."
TX2_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/api/transactions" \
  -H "Content-Type: application/json" \
  -d '{"amount":3000,"category":"еда","description":"ужин","date":"2025-10-21"}')
HTTP_CODE=$(echo "$TX2_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
TX2_BODY=$(echo "$TX2_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $TX2_BODY"
if [ "$HTTP_CODE" != "201" ]; then
    echo "ОШИБКА: Ожидался код 201 (сумма 3450 < 5000)"
    exit 1
fi

echo ""
echo "7.5. Попытка превысить лимит (3450 + 2000 = 5450 > 5000, ожидается 409)..."
TX_EXCEED_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/api/transactions" \
  -H "Content-Type: application/json" \
  -d '{"amount":2000,"category":"еда","description":"превышение","date":"2025-10-22"}')
HTTP_CODE=$(echo "$TX_EXCEED_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
TX_EXCEED_BODY=$(echo "$TX_EXCEED_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $TX_EXCEED_BODY"
if [ "$HTTP_CODE" != "409" ]; then
    echo "ОШИБКА: Ожидался код 409 (превышение бюджета)"
    exit 1
fi

echo ""
echo "7.6. Получение списка транзакций..."
TXS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/api/transactions")
HTTP_CODE=$(echo "$TXS_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
TXS_BODY=$(echo "$TXS_RESPONSE" | sed '/HTTP_CODE/d')
echo "HTTP Code: $HTTP_CODE"
echo "Response: $TXS_BODY"
if [ "$HTTP_CODE" != "200" ]; then
    echo "ОШИБКА: Ожидался код 200"
    exit 1
fi

TX_COUNT=$(echo "$TXS_BODY" | grep -o '"id"' | wc -l | tr -d ' ')
echo "Количество транзакций: $TX_COUNT"
if [ "$TX_COUNT" != "2" ]; then
    echo "ОШИБКА: Ожидалось 2 транзакции, найдено $TX_COUNT"
    exit 1
fi

echo ""
echo "=== Шаг 8: Тестирование кэша отчётов ==="

cat > /tmp/test_report.go <<'EOF'
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/mikhailmogilnikov/go/hw7/ledger"
)

func main() {
	from := time.Date(2025,10,1,0,0,0,0,time.UTC)
	to := time.Date(2025,10,31,0,0,0,0,time.UTC)

	fmt.Println("Первый вызов (cache miss)...")
	start1 := time.Now()
	items1, err := ledger.GetReportSummary(context.Background(), from, to)
	duration1 := time.Since(start1)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Printf("Результат: %+v (время: %v)\n", items1, duration1)

	fmt.Println("\nВторой вызов (cache hit)...")
	start2 := time.Now()
	items2, err := ledger.GetReportSummary(context.Background(), from, to)
	duration2 := time.Since(start2)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Printf("Результат: %+v (время: %v)\n", items2, duration2)

	if duration2 < duration1/2 {
		fmt.Println("\n✓ Кэш работает! Второй запрос быстрее.")
	} else {
		fmt.Println("\n⚠ Кэш может не работать (время похожее)")
	}
}
EOF

cd ledger
go run /tmp/test_report.go
cd ..

echo ""
echo "=== Шаг 9: Проверка персистентности данных ==="
echo "Остановка gateway..."
kill $GATEWAY_PID
wait $GATEWAY_PID 2>/dev/null || true
sleep 1

echo "Перезапуск gateway..."
cd gateway
./gateway > /tmp/gateway.log 2>&1 &
GATEWAY_PID=$!
cd ..
sleep 2

echo "Проверка, что данные сохранились..."
TXS_AFTER_RESTART=$(curl -s "$BASE_URL/api/transactions")
TX_COUNT_AFTER=$(echo "$TXS_AFTER_RESTART" | grep -o '"id"' | wc -l | tr -d ' ')
echo "Количество транзакций после перезапуска: $TX_COUNT_AFTER"
if [ "$TX_COUNT_AFTER" != "2" ]; then
    echo "ОШИБКА: Данные не сохранились после перезапуска"
    exit 1
fi
echo "✓ Данные сохранились после перезапуска!"

echo ""
echo "=== Все тесты пройдены успешно! ==="
echo ""
echo "Логи gateway: /tmp/gateway.log"
echo "Gateway работает на http://localhost:8080"
echo ""
echo "Для остановки gateway выполните: kill $GATEWAY_PID"

