package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	addr := getenv("REDIS_ADDR", "localhost:6379")
	password := os.Getenv("REDIS_PASSWORD")
	db := parseInt(getenv("REDIS_DB", "0"), 0)
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Println("ledger: redis not available:", err)
		return
	}
	log.Println("ledger: connected to redis")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseInt(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}


