package pg

import (
	"context"
	"database/sql"

	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/domain"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Upsert(ctx context.Context, budget domain.Budget) error {
	query := `INSERT INTO budgets(category, limit_amount) 
	          VALUES($1, $2) 
	          ON CONFLICT(category) 
	          DO UPDATE SET limit_amount = EXCLUDED.limit_amount`

	_, err := r.db.ExecContext(ctx, query, budget.Category, budget.Limit)
	return err
}

func (r *BudgetRepository) GetByCategory(ctx context.Context, category string) (*domain.Budget, error) {
	var limit sql.NullFloat64
	err := r.db.QueryRowContext(ctx, `SELECT limit_amount FROM budgets WHERE category=$1`, category).Scan(&limit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if !limit.Valid {
		return nil, nil
	}

	return &domain.Budget{
		Category: category,
		Limit:    limit.Float64,
		Period:   "month",
	}, nil
}

func (r *BudgetRepository) List(ctx context.Context) ([]domain.Budget, error) {
	query := `SELECT category, limit_amount FROM budgets ORDER BY category`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.Budget, 0)
	for rows.Next() {
		var b domain.Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		b.Period = "month"
		result = append(result, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
