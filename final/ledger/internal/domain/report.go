package domain

// CategorySummary сводка по категории
type CategorySummary struct {
	Category         string
	Total            float64
	BudgetLimit      float64
	BudgetPercentage float64 // процент использования бюджета
}



