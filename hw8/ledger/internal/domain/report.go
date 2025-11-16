package domain

// ReportSummary представляет сводку расходов по категориям
type ReportSummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

