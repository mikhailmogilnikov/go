package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/cache"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
)

// LedgerService сервис для работы с финансами
type LedgerService struct {
	txRepo     domain.TransactionRepository
	budgetRepo domain.BudgetRepository
	cache      *cache.Cache
}

// NewLedgerService создаёт новый сервис
func NewLedgerService(txRepo domain.TransactionRepository, budgetRepo domain.BudgetRepository, cache *cache.Cache) *LedgerService {
	return &LedgerService{
		txRepo:     txRepo,
		budgetRepo: budgetRepo,
		cache:      cache,
	}
}

// ErrBudgetExceeded ошибка превышения бюджета
var ErrBudgetExceeded = fmt.Errorf("budget exceeded")

// AddTransaction добавляет транзакцию с проверкой бюджета
// Если бюджет превышен - транзакция отклоняется
func (s *LedgerService) AddTransaction(ctx context.Context, tx *domain.Transaction) (string, error) {
	if err := tx.Validate(); err != nil {
		return "", err
	}

	// Если дата не указана, используем сегодня
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}

	// Проверяем бюджет
	budgetWarning := ""

	budget, err := s.budgetRepo.GetByCategory(ctx, tx.UserID, tx.Category)
	if err != nil {
		return "", err
	}

	if budget != nil {
		// Считаем период для бюджета
		from, to := s.getBudgetPeriod(budget.Period, tx.Date)
		spent, err := s.txRepo.SumByCategory(ctx, tx.UserID, tx.Category, from, to)
		if err != nil {
			return "", err
		}

		newTotal := spent + tx.Amount
		percentage := (newTotal / budget.LimitAmount) * 100

		// Превышение бюджета - отклоняем транзакцию
		if newTotal > budget.LimitAmount {
			return "", fmt.Errorf("%w: limit %.2f, would be %.2f (%.1f%%)",
				ErrBudgetExceeded, budget.LimitAmount, newTotal, percentage)
		}
		
		// Предупреждение если близко к лимиту
		if percentage >= 80 {
			budgetWarning = fmt.Sprintf("Warning: %.1f%% of budget used (%.2f/%.2f)",
				percentage, newTotal, budget.LimitAmount)
		}
	}

	// Создаём транзакцию
	if err := s.txRepo.Create(ctx, tx); err != nil {
		return "", err
	}

	// Инвалидируем кэш отчётов
	if s.cache != nil {
		s.cache.InvalidateReports(ctx, tx.UserID)
	}

	return budgetWarning, nil
}

// GetTransactions возвращает транзакции пользователя
func (s *LedgerService) GetTransactions(ctx context.Context, userID int64, from, to *time.Time, category string) ([]domain.Transaction, error) {
	return s.txRepo.GetByUserID(ctx, userID, from, to, category)
}

// SetBudget устанавливает бюджет
func (s *LedgerService) SetBudget(ctx context.Context, budget *domain.Budget) error {
	if err := budget.Validate(); err != nil {
		return err
	}
	return s.budgetRepo.Upsert(ctx, budget)
}

// GetBudgets возвращает бюджеты пользователя
func (s *LedgerService) GetBudgets(ctx context.Context, userID int64) ([]domain.Budget, error) {
	return s.budgetRepo.GetByUserID(ctx, userID)
}

// GetReport возвращает отчёт по расходам
func (s *LedgerService) GetReport(ctx context.Context, userID int64, from, to time.Time) ([]domain.CategorySummary, float64, error) {
	// Пробуем получить из кэша
	if s.cache != nil {
		if cached, err := s.cache.GetReport(ctx, userID, from, to); err == nil && cached != nil {
			return cached.Categories, cached.TotalExpenses, nil
		}
	}

	// Получаем сводку из БД
	summaries, err := s.txRepo.GetReportSummary(ctx, userID, from, to)
	if err != nil {
		return nil, 0, err
	}

	// Дополняем информацией о бюджетах
	budgets, err := s.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	budgetMap := make(map[string]*domain.Budget)
	for i := range budgets {
		budgetMap[budgets[i].Category] = &budgets[i]
	}

	var totalExpenses float64
	for i := range summaries {
		totalExpenses += summaries[i].Total
		if b, ok := budgetMap[summaries[i].Category]; ok {
			summaries[i].BudgetLimit = b.LimitAmount
			summaries[i].BudgetPercentage = (summaries[i].Total / b.LimitAmount) * 100
		}
	}

	// Сохраняем в кэш
	if s.cache != nil {
		s.cache.SetReport(ctx, userID, from, to, &cache.ReportCache{
			Categories:    summaries,
			TotalExpenses: totalExpenses,
		})
	}

	return summaries, totalExpenses, nil
}

// getBudgetPeriod возвращает начало и конец периода для бюджета
func (s *LedgerService) getBudgetPeriod(period string, date time.Time) (time.Time, time.Time) {
	switch period {
	case "weekly":
		// Начало недели (понедельник)
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from := date.AddDate(0, 0, -(weekday - 1))
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		to := from.AddDate(0, 0, 6)
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())
		return from, to
	default: // monthly
		from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		to := from.AddDate(0, 1, -1)
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())
		return from, to
	}
}

