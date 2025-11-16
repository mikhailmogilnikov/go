package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mikhailmogilnikov/go/hw8/gateway/internal/api"
	"github.com/mikhailmogilnikov/go/hw8/ledger"
)

func main() {
	ctx := context.Background()

	// Инициализация сервиса через фабрику
	service, closeFn, err := ledger.InitService(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}
	defer func() {
		if err := closeFn(); err != nil {
			log.Printf("Error closing resources: %v", err)
		}
	}()

	// Создание обработчиков
	handlers := api.NewHandlers(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTransaction(w, r)
		case http.MethodGet:
			handlers.GetTransactions(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/budgets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateBudget(w, r)
		case http.MethodGet:
			handlers.GetBudgets(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/reports/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetReportSummary(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	handler := api.LoggingMiddleware(mux)

	log.Println("Gateway server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

