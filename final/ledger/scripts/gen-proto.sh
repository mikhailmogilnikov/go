#!/bin/bash

# Генерация Go-кода из proto для Ledger сервиса
# Запускать из директории ledger/

set -e

PROTO_DIR="../api/proto"
OUT_DIR="internal/pb"

mkdir -p "$OUT_DIR/ledger/v1"

protoc \
    --proto_path="$PROTO_DIR/ledger/v1" \
    --go_out="$OUT_DIR/ledger/v1" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR/ledger/v1" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/ledger/v1/ledger.proto"

echo "Ledger proto generation completed!"



