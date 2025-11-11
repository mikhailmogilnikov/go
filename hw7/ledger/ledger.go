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
		return errors.New("сумма транзакции должна быть положительным числом")
	}
	if t.Category == "" {
		return errors.New("категория транзакции не может быть пустой")
	}
	return nil
}

type Budget struct {
	Category string
	Limit    float64
	Period   string
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("лимит бюджета должен быть положительным числом")
	}
	if b.Category == "" {
		return errors.New("категория бюджета не может быть пустой")
	}
	return nil
}

func init() {
	db.Init()
	cache.Init()
}

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	_, err := db.DB.Exec(`INSERT INTO budgets(category, limit_amount) VALUES($1,$2)
ON CONFLICT(category) DO UPDATE SET limit_amount=EXCLUDED.limit_amount`, b.Category, b.Limit)
	if err != nil {
		return err
	}
	if cache.Client != nil {
		cache.Client.Del(context.Background(), "budgets:all")
	}
	return nil
}

func ListBudgets() ([]Budget, error) {
	if cache.Client != nil {
		if v, err := cache.Client.Get(context.Background(), "budgets:all").Result(); err == nil {
			var items []Budget
			if json.Unmarshal([]byte(v), &items) == nil {
				return items, nil
			}
		}
	}
	rows, err := db.DB.Query(`SELECT category, limit_amount FROM budgets ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		b.Period = "месяц"
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if cache.Client != nil {
		if data, err := json.Marshal(result); err == nil {
			cache.Client.Set(context.Background(), "budgets:all", data, 20*time.Second)
		}
	}
	return result, nil
}

func AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}
	var limit sql.NullFloat64
	if err := db.DB.QueryRow(`SELECT limit_amount FROM budgets WHERE category=$1`, tx.Category).Scan(&limit); err != nil && err != sql.ErrNoRows {
		return err
	}
	if limit.Valid {
		var spent float64
		if err := db.DB.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM expenses WHERE category=$1`, tx.Category).Scan(&spent); err != nil {
			return err
		}
		if spent+tx.Amount > limit.Float64 {
			return errors.New("budget exceeded")
		}
	}
	if err := db.DB.QueryRow(`INSERT INTO expenses(amount, category, description, date) VALUES($1,$2,$3,$4) RETURNING id`,
		tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&tx.ID); err != nil {
		return err
	}
	return nil
}

func ListTransactions() ([]Transaction, error) {
	rows, err := db.DB.Query(`SELECT id, amount, category, description, date FROM expenses ORDER BY date DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Category, &t.Description, &t.Date); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

type ReportItem struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func GetReportSummary(ctx context.Context, from, to time.Time) ([]ReportItem, error) {
	key := fmt.Sprintf("report:summary:%s:%s", from.Format("2006-01-02"), to.Format("2006-01-02"))
	if cache.Client != nil {
		if v, err := cache.Client.Get(ctx, key).Result(); err == nil {
			var items []ReportItem
			if json.Unmarshal([]byte(v), &items) == nil {
				return items, nil
			}
		}
	}
	rows, err := db.DB.QueryContext(ctx, `
SELECT category, COALESCE(SUM(amount),0) AS total
FROM expenses
WHERE date BETWEEN $1 AND $2
GROUP BY category
ORDER BY total DESC, category ASC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []ReportItem
	for rows.Next() {
		var it ReportItem
		if err := rows.Scan(&it.Category, &it.Total); err != nil {
			return nil, err
		}
		result = append(result, it)
	}
	if cache.Client != nil {
		if data, err := json.Marshal(result); err == nil {
			cache.Client.Set(ctx, key, data, 30*time.Second)
		}
	}
	return result, nil
}
