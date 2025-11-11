package ledger

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mikhailmogilnikov/go/hw7/ledger/internal/cache"
	"github.com/mikhailmogilnikov/go/hw7/ledger/internal/db"
)

type Validatable interface {
	Validate() error
}

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

func (t Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}
	if t.Category == "" {
		return errors.New("transaction category cannot be empty")
	}
	return nil
}

type Budget struct {
	Category string
	Limit    float64
	Period   string
}

type ReportSummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("budget limit must be positive")
	}
	if b.Category == "" {
		return errors.New("budget category cannot be empty")
	}
	return nil
}

func init() {
	db.Init()
	cache.Init()
}

func AddTransaction(tx *Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	var limitAmount sql.NullFloat64
	err := db.DB().QueryRow(`SELECT limit_amount FROM budgets WHERE category=$1`, tx.Category).Scan(&limitAmount)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil && limitAmount.Valid {
		var spent sql.NullFloat64
		err = db.DB().QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category=$1`, tx.Category).Scan(&spent)
		if err != nil {
			return err
		}

		spentValue := 0.0
		if spent.Valid {
			spentValue = spent.Float64
		}

		if spentValue+tx.Amount > limitAmount.Float64 {
			return errors.New("budget exceeded")
		}
	}

	query := `INSERT INTO expenses(amount, category, description, date) 
	          VALUES($1, $2, $3, $4) 
	          RETURNING id`

	err = db.DB().QueryRow(query, tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&tx.ID)
	if err != nil {
		return err
	}

	return nil
}

func ListTransactions() ([]Transaction, error) {
	query := `SELECT id, amount, category, description, date 
	          FROM expenses 
	          ORDER BY date DESC, id DESC`

	rows, err := db.DB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Transaction, 0)
	for rows.Next() {
		var tx Transaction
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

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO budgets(category, limit_amount) 
	          VALUES($1, $2) 
	          ON CONFLICT(category) 
	          DO UPDATE SET limit_amount = EXCLUDED.limit_amount`

	_, err := db.DB().Exec(query, b.Category, b.Limit)
	if err != nil {
		return err
	}

	client := cache.Client()
	if client != nil {
		ctx := context.Background()
		client.Del(ctx, "budgets:all")
	}

	return nil
}

func ListBudgets() ([]Budget, error) {
	ctx := context.Background()
	cacheKey := "budgets:all"

	client := cache.Client()
	if client != nil {
		cached, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			var result []Budget
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, nil
			}
		}
	}

	query := `SELECT category, limit_amount FROM budgets ORDER BY category`

	rows, err := db.DB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Budget, 0)
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		b.Period = "month"
		result = append(result, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if client != nil {
		data, err := json.Marshal(result)
		if err == nil {
			client.Set(ctx, cacheKey, data, 30*time.Second)
		}
	}

	return result, nil
}

func GetReportSummary(ctx context.Context, from, to time.Time) ([]ReportSummary, error) {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")
	cacheKey := fmt.Sprintf("report:summary:%s:%s", fromStr, toStr)

	client := cache.Client()
	if client != nil {
		cached, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			var result []ReportSummary
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, nil
			}
		}
	}

	query := `SELECT category, COALESCE(SUM(amount), 0) as total 
	          FROM expenses 
	          WHERE date >= $1 AND date <= $2 
	          GROUP BY category 
	          ORDER BY category`

	rows, err := db.DB().QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]ReportSummary, 0)
	for rows.Next() {
		var rs ReportSummary
		if err := rows.Scan(&rs.Category, &rs.Total); err != nil {
			return nil, err
		}
		result = append(result, rs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if client != nil {
		data, err := json.Marshal(result)
		if err == nil {
			client.Set(ctx, cacheKey, data, 30*time.Second)
		}
	}

	return result, nil
}
