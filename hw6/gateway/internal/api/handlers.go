package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mikhailmogilnikov/go/hw6/ledger"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	tx, err := ToTransaction(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err := tx.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err := ledger.AddTransaction(tx); err != nil {
		if strings.Contains(err.Error(), "budget exceeded") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "budget exceeded"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "internal server error"})
		return
	}

	transactions := ledger.ListTransactions()
	if len(transactions) > 0 {
		createdTx := transactions[len(transactions)-1]
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ToTransactionResponse(createdTx))
	}
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	transactions := ledger.ListTransactions()
	response := make([]TransactionResponse, 0, len(transactions))

	for _, tx := range transactions {
		response = append(response, ToTransactionResponse(tx))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateBudget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	budget := ToBudget(req)

	if err := budget.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err := ledger.SetBudget(budget); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToBudgetResponse(budget))
}

func GetBudgets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	budgets := ledger.ListBudgets()
	response := make([]BudgetResponse, 0, len(budgets))

	for _, b := range budgets {
		response = append(response, ToBudgetResponse(b))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
