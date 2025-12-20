#!/bin/bash

# Генерация Go-кода из proto для Gateway
# Запускать из директории gateway/

set -e

PROTO_DIR="../api/proto"
OUT_DIR="internal/pb"

mkdir -p "$OUT_DIR/auth/v1"
mkdir -p "$OUT_DIR/ledger/v1"

# Auth proto
protoc \
    --proto_path="$PROTO_DIR/auth/v1" \
    --go_out="$OUT_DIR/auth/v1" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR/auth/v1" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/auth/v1/auth.proto"

# Ledger proto
protoc \
    --proto_path="$PROTO_DIR/ledger/v1" \
    --go_out="$OUT_DIR/ledger/v1" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR/ledger/v1" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/ledger/v1/ledger.proto"

echo "Gateway proto generation completed!"



