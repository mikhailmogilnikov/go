package domain

import (
	"context"
	"time"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByUserID(ctx context.Context, userID int64, from, to *time.Time, category string) ([]Transaction, error)
	SumByCategory(ctx context.Context, userID int64, category string, from, to time.Time) (float64, error)
	GetReportSummary(ctx context.Context, userID int64, from, to time.Time) ([]CategorySummary, error)
}

type BudgetRepository interface {
	Upsert(ctx context.Context, budget *Budget) error
	GetByUserID(ctx context.Context, userID int64) ([]Budget, error)
	GetByCategory(ctx context.Context, userID int64, category string) (*Budget, error)
}



