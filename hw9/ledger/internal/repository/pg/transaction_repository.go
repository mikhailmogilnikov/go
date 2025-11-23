package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx domain.Transaction) (int, error) {
	query := `INSERT INTO expenses(amount, category, description, date) 
	          VALUES($1, $2, $3, $4) 
	          RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TransactionRepository) List(ctx context.Context) ([]domain.Transaction, error) {
	query := `SELECT id, amount, category, description, date 
	          FROM expenses 
	          ORDER BY date DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.Transaction, 0)
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Category, &tx.Description, &tx.Date); err != nil {
			return nil, err
		}
		result = append(result, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TransactionRepository) SumByCategory(ctx context.Context, category string) (float64, error) {
	var spent sql.NullFloat64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category=$1`, category).Scan(&spent)
	if err != nil {
		return 0, err
	}

	if !spent.Valid {
		return 0, nil
	}

	return spent.Float64, nil
}

func (r *TransactionRepository) SumByCategoryPeriod(ctx context.Context, category string, from, to time.Time) (float64, error) {
	var spent sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category=$1 AND date >= $2 AND date <= $3`,
		category, from, to).Scan(&spent)
	if err != nil {
		return 0, err
	}

	if !spent.Valid {
		return 0, nil
	}

	return spent.Float64, nil
}

func (r *TransactionRepository) GetReportSummary(ctx context.Context, from, to time.Time) ([]domain.ReportSummary, error) {
	query := `SELECT category, COALESCE(SUM(amount), 0) as total 
	          FROM expenses 
	          WHERE date >= $1 AND date <= $2 
	          GROUP BY category 
	          ORDER BY category`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.ReportSummary, 0)
	for rows.Next() {
		var rs domain.ReportSummary
		if err := rows.Scan(&rs.Category, &rs.Total); err != nil {
			return nil, err
		}
		result = append(result, rs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TransactionRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT category FROM expenses
		UNION
		SELECT DISTINCT category FROM budgets
		ORDER BY category
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]string, 0)
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
