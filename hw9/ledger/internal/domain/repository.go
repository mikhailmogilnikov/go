package domain

import (
	"context"
	"time"
)

type BudgetRepository interface {
	Upsert(ctx context.Context, budget Budget) error
	GetByCategory(ctx context.Context, category string) (*Budget, error)
	List(ctx context.Context) ([]Budget, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx Transaction) (int, error)
	List(ctx context.Context) ([]Transaction, error)
	SumByCategory(ctx context.Context, category string) (float64, error)
	SumByCategoryPeriod(ctx context.Context, category string, from, to time.Time) (float64, error)
	GetReportSummary(ctx context.Context, from, to time.Time) ([]ReportSummary, error)
	GetCategories(ctx context.Context) ([]string, error)
}

