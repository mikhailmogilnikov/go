package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestRegisterRequest_Validation проверяет валидацию запроса регистрации
func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid request",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusOK, // валидация прошла
		},
		{
			name: "missing email",
			body: map[string]interface{}{
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			body: map[string]interface{}{
				"email":    "not-an-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			// Простой хендлер для тестирования валидации
			router.POST("/register", func(c *gin.Context) {
				var req RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"email": req.Email})
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestAddTransactionRequest_Validation проверяет валидацию запроса транзакции
func TestAddTransactionRequest_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid request",
			body: map[string]interface{}{
				"amount":   1500.50,
				"category": "food",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "with description and date",
			body: map[string]interface{}{
				"amount":      1500.50,
				"category":    "food",
				"description": "Lunch",
				"date":        "2024-12-15",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing amount",
			body: map[string]interface{}{
				"category": "food",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "zero amount",
			body: map[string]interface{}{
				"amount":   0,
				"category": "food",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative amount",
			body: map[string]interface{}{
				"amount":   -100,
				"category": "food",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing category",
			body: map[string]interface{}{
				"amount": 1500,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/transaction", func(c *gin.Context) {
				var req AddTransactionRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"amount": req.Amount, "category": req.Category})
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

// TestSetBudgetRequest_Validation проверяет валидацию запроса бюджета
func TestSetBudgetRequest_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid monthly budget",
			body: map[string]interface{}{
				"category":     "food",
				"limit_amount": 15000,
				"period":       "monthly",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "valid weekly budget",
			body: map[string]interface{}{
				"category":     "transport",
				"limit_amount": 5000,
				"period":       "weekly",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing category",
			body: map[string]interface{}{
				"limit_amount": 15000,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing limit",
			body: map[string]interface{}{
				"category": "food",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "zero limit",
			body: map[string]interface{}{
				"category":     "food",
				"limit_amount": 0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/budget", func(c *gin.Context) {
				var req SetBudgetRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"category": req.Category, "limit": req.LimitAmount})
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/budget", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

