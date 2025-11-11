package main

import (
	"log"
	"net/http"

	"github.com/mikhailmogilnikov/go/hw7/gateway/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.CreateTransaction(w, r)
		case http.MethodGet:
			api.GetTransactions(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/budgets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.CreateBudget(w, r)
		case http.MethodGet:
			api.GetBudgets(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/reports/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.GetReportSummary(w, r)
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
