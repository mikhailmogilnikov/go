package ledger

import (
	"errors"
	"time"
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

var (
	transactions []Transaction
	budgets      map[string]Budget
)

func init() {
	budgets = make(map[string]Budget)
}

func getCategoryTotal(category string) float64 {
	var total float64
	for _, tx := range transactions {
		if tx.Category == category {
			total += tx.Amount
		}
	}
	return total
}

func AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	if budget, exists := budgets[tx.Category]; exists {
		currentTotal := getCategoryTotal(tx.Category)
		newTotal := currentTotal + tx.Amount

		if newTotal > budget.Limit {
			return errors.New("budget exceeded")
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	budgets[b.Category] = b
	return nil
}

func ListBudgets() []Budget {
	result := make([]Budget, 0, len(budgets))
	for _, budget := range budgets {
		result = append(result, budget)
	}
	return result
}

