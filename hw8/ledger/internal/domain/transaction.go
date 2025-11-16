package domain

import (
	"errors"
	"time"
)

// Transaction представляет транзакцию расходов
type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

// Validate проверяет валидность транзакции
func (t Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}
	if t.Category == "" {
		return errors.New("transaction category cannot be empty")
	}
	return nil
}

