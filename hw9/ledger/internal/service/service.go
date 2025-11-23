package service

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/domain"
)

type Service interface {
	SetBudget(ctx context.Context, budget domain.Budget) error
	GetBudgets(ctx context.Context) ([]domain.Budget, error)
	AddTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	GetTransactions(ctx context.Context) ([]domain.Transaction, error)
	GetReportSummary(ctx context.Context, from, to time.Time) (map[string]float64, error)
}

type LedgerService struct {
	budgetRepo     domain.BudgetRepository
	transactionRepo domain.TransactionRepository
}

func NewLedgerService(budgetRepo domain.BudgetRepository, transactionRepo domain.TransactionRepository) *LedgerService {
	return &LedgerService{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *LedgerService) SetBudget(ctx context.Context, budget domain.Budget) error {
	if err := budget.Validate(); err != nil {
		return err
	}

	return s.budgetRepo.Upsert(ctx, budget)
}

func (s *LedgerService) GetBudgets(ctx context.Context) ([]domain.Budget, error) {
	return s.budgetRepo.List(ctx)
}

func (s *LedgerService) AddTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	if err := tx.Validate(); err != nil {
		return tx, err
	}

	budget, err := s.budgetRepo.GetByCategory(ctx, tx.Category)
	if err != nil {
		return tx, err
	}

	if budget != nil {
		spent, err := s.transactionRepo.SumByCategory(ctx, tx.Category)
		if err != nil {
			return tx, err
		}

		if spent+tx.Amount > budget.Limit {
			return tx, errors.New("budget exceeded")
		}
	}

	id, err := s.transactionRepo.Create(ctx, tx)
	if err != nil {
		return tx, err
	}

	tx.ID = id
	return tx, nil
}

func (s *LedgerService) GetTransactions(ctx context.Context) ([]domain.Transaction, error) {
	return s.transactionRepo.List(ctx)
}

func (s *LedgerService) GetReportSummary(ctx context.Context, from, to time.Time) (map[string]float64, error) {
	categories, err := s.transactionRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return make(map[string]float64), nil
	}

	heartbeatDone := make(chan bool)
	ticker := time.NewTicker(400 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Calculating report summary...")
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-heartbeatDone:
				ticker.Stop()
				return
			}
		}
	}()

	type result struct {
		category string
		total    float64
		err      error
	}

	results := make(chan result, len(categories))
	var wg sync.WaitGroup

	for _, category := range categories {
		wg.Add(1)
		go func(cat string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results <- result{category: cat, err: ctx.Err()}
				return
			default:
			}

			total, err := s.transactionRepo.SumByCategoryPeriod(ctx, cat, from, to)
			results <- result{category: cat, total: total, err: err}
		}(category)
	}

	go func() {
		wg.Wait()
		close(results)
		close(heartbeatDone)
	}()

	report := make(map[string]float64)
	for res := range results {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if res.err != nil {
			return nil, res.err
		}

		report[res.category] = res.total
	}

	return report, nil
}

