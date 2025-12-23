package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mikhailmogilnikov/go/final/gateway/internal/config"
	"github.com/mikhailmogilnikov/go/final/gateway/internal/handler"
	"github.com/mikhailmogilnikov/go/final/gateway/internal/middleware"
	authv1 "github.com/mikhailmogilnikov/go/final/gateway/internal/pb/auth/v1"
	ledgerv1 "github.com/mikhailmogilnikov/go/final/gateway/internal/pb/ledger/v1"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("Connecting to Auth gRPC at %s", cfg.AuthGRPCAddr)
	authConn, err := grpc.DialContext(ctx, cfg.AuthGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Auth service: %v", err)
	}
	defer authConn.Close()
	log.Println("Connected to Auth service")

	log.Printf("Connecting to Ledger gRPC at %s", cfg.LedgerGRPCAddr)
	ledgerConn, err := grpc.DialContext(ctx, cfg.LedgerGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Ledger service: %v", err)
	}
	defer ledgerConn.Close()
	log.Println("Connected to Ledger service")

	authClient := authv1.NewAuthServiceClient(authConn)
	ledgerClient := ledgerv1.NewLedgerServiceClient(ledgerConn)

	authMiddleware := middleware.NewAuthMiddleware(authClient)
	authHandler := handler.NewAuthHandler(authClient)
	ledgerHandler := handler.NewLedgerHandler(ledgerClient)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := engine.Group("/api")
	authHandler.RegisterRoutes(api)
	ledgerHandler.RegisterRoutes(api, authMiddleware)

	httpAddr := ":" + cfg.HTTPPort
	server := &http.Server{
		Addr:         httpAddr,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Printf("Gateway HTTP server listening on %s", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Gateway...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	log.Println("Gateway stopped")
}

