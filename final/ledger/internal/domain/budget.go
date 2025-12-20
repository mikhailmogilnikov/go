package domain

import "errors"

// Budget представляет бюджет по категории
type Budget struct {
	ID          int64
	UserID      int64
	Category    string
	LimitAmount float64
	Period      string // "monthly" или "weekly"
}

// Validate проверяет валидность бюджета
func (b *Budget) Validate() error {
	if b.LimitAmount <= 0 {
		return errors.New("limit must be positive")
	}
	if b.Category == "" {
		return errors.New("category is required")
	}
	if b.UserID <= 0 {
		return errors.New("user_id is required")
	}
	if b.Period != "monthly" && b.Period != "weekly" {
		b.Period = "monthly" // дефолт
	}
	return nil
}



