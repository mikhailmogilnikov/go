package handler

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mikhailmogilnikov/go/final/gateway/internal/middleware"
	ledgerv1 "github.com/mikhailmogilnikov/go/final/gateway/internal/pb/ledger/v1"
)

// LedgerHandler хендлер для работы с финансами
type LedgerHandler struct {
	ledgerClient ledgerv1.LedgerServiceClient
}

// NewLedgerHandler создаёт новый хендлер
func NewLedgerHandler(ledgerClient ledgerv1.LedgerServiceClient) *LedgerHandler {
	return &LedgerHandler{ledgerClient: ledgerClient}
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// === Транзакции ===

// AddTransactionRequest запрос на добавление транзакции
type AddTransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
	Date        string  `json:"date"` // формат YYYY-MM-DD
}

// TransactionResponse ответ с транзакцией
type TransactionResponse struct {
	ID            int64   `json:"id"`
	Amount        float64 `json:"amount"`
	Category      string  `json:"category"`
	Description   string  `json:"description"`
	Date          string  `json:"date"`
	BudgetWarning string  `json:"budget_warning,omitempty"`
}

// AddTransaction добавляет транзакцию
// @Summary Добавить транзакцию
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddTransactionRequest true "Данные транзакции"
// @Success 201 {object} TransactionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/transactions [post]
func (h *LedgerHandler) AddTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим дату
	var date *timestamppb.Timestamp
	if req.Date != "" {
		t, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		date = timestamppb.New(t)
	}

	resp, err := h.ledgerClient.AddTransaction(c.Request.Context(), &ledgerv1.AddTransactionRequest{
		UserId:      userID,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        date,
	})
	if err != nil {
		// Превышение бюджета - возвращаем 409 Conflict
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := resp.GetTransaction()
	c.JSON(http.StatusCreated, TransactionResponse{
		ID:             tx.GetId(),
		Amount:         tx.GetAmount(),
		Category:       tx.GetCategory(),
		Description:    tx.GetDescription(),
		Date:           tx.GetDate().AsTime().Format("2006-01-02"),
		BudgetWarning:  resp.GetBudgetWarning(),
	})
}

// GetTransactions возвращает транзакции
// @Summary Получить транзакции
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param from query string false "Дата начала (YYYY-MM-DD)"
// @Param to query string false "Дата конца (YYYY-MM-DD)"
// @Param category query string false "Фильтр по категории"
// @Success 200 {array} TransactionResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/transactions [get]
func (h *LedgerHandler) GetTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := &ledgerv1.GetTransactionsRequest{
		UserId:   userID,
		Category: c.Query("category"),
	}

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err == nil {
			req.From = timestamppb.New(t)
		}
	}
	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err == nil {
			req.To = timestamppb.New(t)
		}
	}

	resp, err := h.ledgerClient.GetTransactions(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactions := make([]TransactionResponse, 0, len(resp.GetTransactions()))
	for _, tx := range resp.GetTransactions() {
		transactions = append(transactions, TransactionResponse{
			ID:          tx.GetId(),
			Amount:      tx.GetAmount(),
			Category:    tx.GetCategory(),
			Description: tx.GetDescription(),
			Date:        tx.GetDate().AsTime().Format("2006-01-02"),
		})
	}

	c.JSON(http.StatusOK, transactions)
}

// === Бюджеты ===

// SetBudgetRequest запрос на установку бюджета
type SetBudgetRequest struct {
	Category    string  `json:"category" binding:"required"`
	LimitAmount float64 `json:"limit_amount" binding:"required,gt=0"`
	Period      string  `json:"period"` // monthly или weekly
}

// BudgetResponse ответ с бюджетом
type BudgetResponse struct {
	ID          int64   `json:"id"`
	Category    string  `json:"category"`
	LimitAmount float64 `json:"limit_amount"`
	Period      string  `json:"period"`
}

// SetBudget устанавливает бюджет
// @Summary Установить бюджет
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetBudgetRequest true "Данные бюджета"
// @Success 201 {object} BudgetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/budgets [post]
func (h *LedgerHandler) SetBudget(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SetBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	resp, err := h.ledgerClient.SetBudget(c.Request.Context(), &ledgerv1.SetBudgetRequest{
		UserId:      userID,
		Category:    req.Category,
		LimitAmount: req.LimitAmount,
		Period:      req.Period,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budget := resp.GetBudget()
	c.JSON(http.StatusCreated, BudgetResponse{
		ID:          budget.GetId(),
		Category:    budget.GetCategory(),
		LimitAmount: budget.GetLimitAmount(),
		Period:      budget.GetPeriod(),
	})
}

// GetBudgets возвращает бюджеты
// @Summary Получить бюджеты
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Success 200 {array} BudgetResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/budgets [get]
func (h *LedgerHandler) GetBudgets(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.ledgerClient.GetBudgets(c.Request.Context(), &ledgerv1.GetBudgetsRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budgets := make([]BudgetResponse, 0, len(resp.GetBudgets()))
	for _, b := range resp.GetBudgets() {
		budgets = append(budgets, BudgetResponse{
			ID:          b.GetId(),
			Category:    b.GetCategory(),
			LimitAmount: b.GetLimitAmount(),
			Period:      b.GetPeriod(),
		})
	}

	c.JSON(http.StatusOK, budgets)
}

// === Отчёты ===

// CategorySummaryResponse сводка по категории
type CategorySummaryResponse struct {
	Category         string  `json:"category"`
	Total            float64 `json:"total"`
	BudgetLimit      float64 `json:"budget_limit,omitempty"`
	BudgetPercentage float64 `json:"budget_percentage,omitempty"`
}

// ReportResponse ответ с отчётом
type ReportResponse struct {
	Categories    []CategorySummaryResponse `json:"categories"`
	TotalExpenses float64                   `json:"total_expenses"`
	From          string                    `json:"from"`
	To            string                    `json:"to"`
}

// GetReport возвращает отчёт
// @Summary Получить отчёт по расходам
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string true "Дата начала (YYYY-MM-DD)"
// @Param to query string true "Дата конца (YYYY-MM-DD)"
// @Success 200 {object} ReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/reports [get]
func (h *LedgerHandler) GetReport(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to parameters are required (YYYY-MM-DD)"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format"})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format"})
		return
	}

	resp, err := h.ledgerClient.GetReport(c.Request.Context(), &ledgerv1.GetReportRequest{
		UserId: userID,
		From:   timestamppb.New(from),
		To:     timestamppb.New(to),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories := make([]CategorySummaryResponse, 0, len(resp.GetCategories()))
	for _, cat := range resp.GetCategories() {
		categories = append(categories, CategorySummaryResponse{
			Category:         cat.GetCategory(),
			Total:            cat.GetTotal(),
			BudgetLimit:      cat.GetBudgetLimit(),
			BudgetPercentage: cat.GetBudgetPercentage(),
		})
	}

	c.JSON(http.StatusOK, ReportResponse{
		Categories:    categories,
		TotalExpenses: resp.GetTotalExpenses(),
		From:          fromStr,
		To:            toStr,
	})
}

// === CSV ===

// ImportCSVRequest запрос на импорт CSV
type ImportCSVRequest struct {
	CSVData string `json:"csv_data" binding:"required"` // base64 encoded
}

// ImportCSVResponse ответ на импорт CSV
type ImportCSVResponse struct {
	ImportedCount int32    `json:"imported_count"`
	SkippedCount  int32    `json:"skipped_count"`
	Errors        []string `json:"errors,omitempty"`
}

// ImportCSV импортирует транзакции из CSV
// @Summary Импорт транзакций из CSV
// @Tags csv
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ImportCSVRequest true "CSV данные в base64"
// @Success 200 {object} ImportCSVResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/csv/import [post]
func (h *LedgerHandler) ImportCSV(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ImportCSVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Декодируем base64
	csvData, err := base64.StdEncoding.DecodeString(req.CSVData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 data"})
		return
	}

	resp, err := h.ledgerClient.ImportCSV(c.Request.Context(), &ledgerv1.ImportCSVRequest{
		UserId:  userID,
		CsvData: csvData,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ImportCSVResponse{
		ImportedCount: resp.GetImportedCount(),
		SkippedCount:  resp.GetSkippedCount(),
		Errors:        resp.GetErrors(),
	})
}

// ExportCSV экспортирует транзакции в CSV
// @Summary Экспорт транзакций в CSV
// @Tags csv
// @Produce json
// @Security BearerAuth
// @Param from query string false "Дата начала (YYYY-MM-DD)"
// @Param to query string false "Дата конца (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /api/csv/export [get]
func (h *LedgerHandler) ExportCSV(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := &ledgerv1.ExportCSVRequest{
		UserId: userID,
	}

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err == nil {
			req.From = timestamppb.New(t)
		}
	}
	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err == nil {
			req.To = timestamppb.New(t)
		}
	}

	resp, err := h.ledgerClient.ExportCSV(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"csv_data":   base64.StdEncoding.EncodeToString(resp.GetCsvData()),
		"rows_count": resp.GetRowsCount(),
	})
}

// RegisterRoutes регистрирует роуты
func (h *LedgerHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	// Транзакции
	transactions := r.Group("/transactions")
	transactions.Use(authMiddleware.RequireAuth())
	{
		transactions.POST("", h.AddTransaction)
		transactions.GET("", h.GetTransactions)
	}

	// Бюджеты
	budgets := r.Group("/budgets")
	budgets.Use(authMiddleware.RequireAuth())
	{
		budgets.POST("", h.SetBudget)
		budgets.GET("", h.GetBudgets)
	}

	// Отчёты
	reports := r.Group("/reports")
	reports.Use(authMiddleware.RequireAuth())
	{
		reports.GET("", h.GetReport)
	}

	// CSV
	csv := r.Group("/csv")
	csv.Use(authMiddleware.RequireAuth())
	{
		csv.POST("/import", h.ImportCSV)
		csv.GET("/export", h.ExportCSV)
	}
}

