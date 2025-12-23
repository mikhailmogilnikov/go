package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		tx.UserID, tx.Amount, tx.Category, tx.Description, tx.Date,
	).Scan(&tx.ID, &tx.CreatedAt)
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userID int64, from, to *time.Time, category string) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, category, description, date, created_at
		FROM transactions
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argNum := 2

	if from != nil {
		query += ` AND date >= $` + string(rune('0'+argNum))
		args = append(args, *from)
		argNum++
	}
	if to != nil {
		query += ` AND date <= $` + string(rune('0'+argNum))
		args = append(args, *to)
		argNum++
	}
	if category != "" {
		query += ` AND category = $` + string(rune('0'+argNum))
		args = append(args, category)
	}
	query += ` ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Category, &tx.Description, &tx.Date, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) SumByCategory(ctx context.Context, userID int64, category string, from, to time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1 AND category = $2 AND date >= $3 AND date <= $4
	`
	var sum float64
	err := r.db.QueryRow(ctx, query, userID, category, from, to).Scan(&sum)
	return sum, err
}

func (r *TransactionRepository) GetReportSummary(ctx context.Context, userID int64, from, to time.Time) ([]domain.CategorySummary, error) {
	query := `
		SELECT category, SUM(amount) as total
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY category
		ORDER BY total DESC
	`
	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []domain.CategorySummary
	for rows.Next() {
		var s domain.CategorySummary
		err := rows.Scan(&s.Category, &s.Total)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}



