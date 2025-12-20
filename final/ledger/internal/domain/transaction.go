package domain

import (
	"errors"
	"time"
)

// Transaction представляет транзакцию расходов
type Transaction struct {
	ID          int64
	UserID      int64
	Amount      float64
	Category    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

// Validate проверяет валидность транзакции
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if t.Category == "" {
		return errors.New("category is required")
	}
	if t.UserID <= 0 {
		return errors.New("user_id is required")
	}
	return nil
}



