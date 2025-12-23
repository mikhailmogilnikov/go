package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/cache"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/config"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/grpcserver"
	pb "github.com/mikhailmogilnikov/go/final/ledger/internal/pb/ledger/v1"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/repository/pg"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	var redisCache *cache.Cache
	redisCache, err = cache.NewCache(cfg.RedisAddr, cfg.RedisDB, cfg.RedisPassword, cfg.CacheTTL)
	if err != nil {
		log.Printf("Warning: Redis not available, caching disabled: %v", err)
		redisCache = nil
	} else {
		defer redisCache.Close()
	}

	txRepo := pg.NewTransactionRepository(pool)
	budgetRepo := pg.NewBudgetRepository(pool)
	ledgerService := service.NewLedgerService(txRepo, budgetRepo, redisCache)
	ledgerServer := grpcserver.NewLedgerServer(ledgerService)

	grpcAddr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterLedgerServiceServer(grpcSrv, ledgerServer)

	go func() {
		log.Printf("Ledger gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Ledger service...")
	grpcSrv.GracefulStop()
	log.Println("Ledger service stopped")
}
