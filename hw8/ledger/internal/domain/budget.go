package domain

import "errors"

// Budget представляет бюджет по категории
type Budget struct {
	Category string
	Limit    float64
	Period   string
}

// Validate проверяет валидность бюджета
func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("budget limit must be positive")
	}
	if b.Category == "" {
		return errors.New("budget category cannot be empty")
	}
	return nil
}

