package domain

import (
	"context"
	"time"
)

// BudgetRepository определяет интерфейс для работы с бюджетами
type BudgetRepository interface {
	// Upsert создает или обновляет бюджет по категории
	Upsert(ctx context.Context, budget Budget) error
	// GetByCategory возвращает бюджет по категории
	GetByCategory(ctx context.Context, category string) (*Budget, error)
	// List возвращает все бюджеты
	List(ctx context.Context) ([]Budget, error)
}

// TransactionRepository определяет интерфейс для работы с транзакциями
type TransactionRepository interface {
	// Create создает новую транзакцию и возвращает её ID
	Create(ctx context.Context, tx Transaction) (int, error)
	// List возвращает все транзакции
	List(ctx context.Context) ([]Transaction, error)
	// SumByCategory возвращает сумму расходов по категории
	SumByCategory(ctx context.Context, category string) (float64, error)
	// GetReportSummary возвращает сводку расходов за период
	GetReportSummary(ctx context.Context, from, to time.Time) ([]ReportSummary, error)
}

