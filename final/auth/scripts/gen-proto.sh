#!/bin/bash

# Генерация Go-кода из proto для Auth сервиса
# Запускать из директории auth/

set -e

PROTO_DIR="../api/proto"
OUT_DIR="internal/pb"

mkdir -p "$OUT_DIR/auth/v1"

protoc \
    --proto_path="$PROTO_DIR/auth/v1" \
    --go_out="$OUT_DIR/auth/v1" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR/auth/v1" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/auth/v1/auth.proto"

echo "Auth proto generation completed!"



