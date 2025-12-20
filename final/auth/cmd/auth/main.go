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

	"github.com/mikhailmogilnikov/go/final/auth/internal/config"
	"github.com/mikhailmogilnikov/go/final/auth/internal/grpcserver"
	pb "github.com/mikhailmogilnikov/go/final/auth/internal/pb/auth/v1"
	"github.com/mikhailmogilnikov/go/final/auth/internal/repository/pg"
	"github.com/mikhailmogilnikov/go/final/auth/internal/service"
)

func main() {
	// Загружаем .env если есть
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Подключаемся к БД
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Создаём зависимости
	userRepo := pg.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenTTL)
	authServer := grpcserver.NewAuthServer(authService)

	// Запускаем gRPC сервер
	grpcAddr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcSrv, authServer)

	// Запускаем в горутине
	go func() {
		log.Printf("Auth gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Ждём сигнала завершения
	<-ctx.Done()
	log.Println("Shutting down Auth service...")
	grpcSrv.GracefulStop()
	log.Println("Auth service stopped")
}



