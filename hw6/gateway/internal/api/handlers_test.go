package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mikhailmogilnikov/go/hw6/ledger"
)

func setupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateTransaction(w, r)
		case http.MethodGet:
			GetTransactions(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/budgets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateBudget(w, r)
		case http.MethodGet:
			GetBudgets(w, r)
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return LoggingMiddleware(mux)
}

func TestCreateBudget(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		reqBody := CreateBudgetRequest{
			Category: "еда",
			Limit:    5000.0,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("CreateBudget() status = %d, want %d", w.Code, http.StatusCreated)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json; charset=utf-8" {
			t.Errorf("CreateBudget() Content-Type = %s, want 'application/json; charset=utf-8'", contentType)
		}

		var resp BudgetResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateBudget() decode error = %v", err)
		}

		if resp.Category != "еда" {
			t.Errorf("CreateBudget() Category = %s, want 'еда'", resp.Category)
		}
		if resp.Limit != 5000.0 {
			t.Errorf("CreateBudget() Limit = %f, want 5000.0", resp.Limit)
		}
	})

	t.Run("invalid_limit", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		reqBody := CreateBudgetRequest{
			Category: "еда",
			Limit:    -100.0,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("CreateBudget() status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateBudget() decode error = %v", err)
		}

		if resp.Error == "" {
			t.Error("CreateBudget() error message is empty")
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("CreateBudget() status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateBudget() decode error = %v", err)
		}

		if resp.Error != "invalid request body" {
			t.Errorf("CreateBudget() error = %s, want 'invalid request body'", resp.Error)
		}
	})

	t.Run("get_budgets", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		budget := ledger.Budget{Category: "еда", Limit: 5000.0}
		ledger.SetBudget(budget)

		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetBudgets() status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp []BudgetResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("GetBudgets() decode error = %v", err)
		}

		if len(resp) != 1 {
			t.Errorf("GetBudgets() len = %d, want 1", len(resp))
		}

		if resp[0].Category != "еда" {
			t.Errorf("GetBudgets() Category = %s, want 'еда'", resp[0].Category)
		}
	})
}

func TestCreateTransaction(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		budget := ledger.Budget{Category: "еда", Limit: 5000.0}
		ledger.SetBudget(budget)

		router := setupRouter()

		reqBody := CreateTransactionRequest{
			Amount:      450.0,
			Category:    "еда",
			Description: "ланч",
			Date:        "2025-10-20",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("CreateTransaction() status = %d, want %d", w.Code, http.StatusCreated)
		}

		var resp TransactionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateTransaction() decode error = %v", err)
		}

		if resp.Amount != 450.0 {
			t.Errorf("CreateTransaction() Amount = %f, want 450.0", resp.Amount)
		}
		if resp.Category != "еда" {
			t.Errorf("CreateTransaction() Category = %s, want 'еда'", resp.Category)
		}
		if resp.Date != "2025-10-20" {
			t.Errorf("CreateTransaction() Date = %s, want '2025-10-20'", resp.Date)
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		budget := ledger.Budget{Category: "еда", Limit: 5000.0}
		ledger.SetBudget(budget)

		tx1 := ledger.Transaction{
			Amount:   4800.0,
			Category: "еда",
			Date:     time.Now(),
		}
		ledger.AddTransaction(tx1)

		router := setupRouter()

		reqBody := CreateTransactionRequest{
			Amount:      300.0,
			Category:    "еда",
			Description: "перекус",
			Date:        "2025-10-20",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("CreateTransaction() status = %d, want %d", w.Code, http.StatusConflict)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateTransaction() decode error = %v", err)
		}

		if resp.Error != "budget exceeded" {
			t.Errorf("CreateTransaction() error = %s, want 'budget exceeded'", resp.Error)
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("CreateTransaction() status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateTransaction() decode error = %v", err)
		}

		if resp.Error != "invalid request body" {
			t.Errorf("CreateTransaction() error = %s, want 'invalid request body'", resp.Error)
		}
	})

	t.Run("invalid_amount", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		reqBody := CreateTransactionRequest{
			Amount:      -100.0,
			Category:    "еда",
			Description: "ланч",
			Date:        "2025-10-20",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("CreateTransaction() status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("CreateTransaction() decode error = %v", err)
		}

		if resp.Error == "" {
			t.Error("CreateTransaction() error message is empty")
		}
	})

	t.Run("get_transactions", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		tx := ledger.Transaction{
			Amount:      450.0,
			Category:    "еда",
			Description: "ланч",
			Date:        time.Now(),
		}
		ledger.AddTransaction(tx)

		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetTransactions() status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp []TransactionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("GetTransactions() decode error = %v", err)
		}

		if len(resp) != 1 {
			t.Errorf("GetTransactions() len = %d, want 1", len(resp))
		}

		if resp[0].Amount != 450.0 {
			t.Errorf("GetTransactions() Amount = %f, want 450.0", resp[0].Amount)
		}
	})
}

func TestIntegrationFlow(t *testing.T) {
	t.Run("create_budget_and_transaction", func(t *testing.T) {
		ledger.Reset()
		defer ledger.Reset()

		router := setupRouter()

		budgetReq := CreateBudgetRequest{
			Category: "еда",
			Limit:    5000.0,
		}
		body, _ := json.Marshal(budgetReq)
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("CreateBudget() status = %d, want %d", w.Code, http.StatusCreated)
		}

		txReq := CreateTransactionRequest{
			Amount:      450.0,
			Category:    "еда",
			Description: "ланч",
			Date:        "2025-10-20",
		}
		body, _ = json.Marshal(txReq)
		req = httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("CreateTransaction() status = %d, want %d", w.Code, http.StatusCreated)
		}

		req = httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GetTransactions() status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp []TransactionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("GetTransactions() decode error = %v", err)
		}

		if len(resp) != 1 {
			t.Errorf("GetTransactions() len = %d, want 1", len(resp))
		}
	})
}
