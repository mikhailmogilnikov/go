package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
)

type BudgetRepository struct {
	db *pgxpool.Pool
}

func NewBudgetRepository(db *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Upsert(ctx context.Context, budget *domain.Budget) error {
	query := `
		INSERT INTO budgets (user_id, category, limit_amount, period)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, category) DO UPDATE SET
			limit_amount = EXCLUDED.limit_amount,
			period = EXCLUDED.period
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		budget.UserID, budget.Category, budget.LimitAmount, budget.Period,
	).Scan(&budget.ID)
}

func (r *BudgetRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.Budget, error) {
	query := `
		SELECT id, user_id, category, limit_amount, period
		FROM budgets
		WHERE user_id = $1
		ORDER BY category
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []domain.Budget
	for rows.Next() {
		var b domain.Budget
		err := rows.Scan(&b.ID, &b.UserID, &b.Category, &b.LimitAmount, &b.Period)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, rows.Err()
}

func (r *BudgetRepository) GetByCategory(ctx context.Context, userID int64, category string) (*domain.Budget, error) {
	query := `
		SELECT id, user_id, category, limit_amount, period
		FROM budgets
		WHERE user_id = $1 AND category = $2
	`
	var b domain.Budget
	err := r.db.QueryRow(ctx, query, userID, category).
		Scan(&b.ID, &b.UserID, &b.Category, &b.LimitAmount, &b.Period)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}



