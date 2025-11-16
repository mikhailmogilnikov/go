package service

import (
	"context"
	"errors"
	"time"

	"github.com/mikhailmogilnikov/go/hw8/ledger/internal/domain"
)

// Service определяет интерфейс приложения для Gateway
type Service interface {
	// SetBudget устанавливает или обновляет бюджет
	SetBudget(ctx context.Context, budget domain.Budget) error
	// GetBudgets возвращает все бюджеты
	GetBudgets(ctx context.Context) ([]domain.Budget, error)
	// AddTransaction добавляет новую транзакцию
	AddTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	// GetTransactions возвращает все транзакции
	GetTransactions(ctx context.Context) ([]domain.Transaction, error)
	// GetReportSummary возвращает сводку расходов за период
	GetReportSummary(ctx context.Context, from, to time.Time) ([]domain.ReportSummary, error)
}

// LedgerService реализует Service с бизнес-логикой
type LedgerService struct {
	budgetRepo     domain.BudgetRepository
	transactionRepo domain.TransactionRepository
}

// NewLedgerService создает новый экземпляр LedgerService
func NewLedgerService(budgetRepo domain.BudgetRepository, transactionRepo domain.TransactionRepository) *LedgerService {
	return &LedgerService{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
	}
}

// SetBudget устанавливает или обновляет бюджет
func (s *LedgerService) SetBudget(ctx context.Context, budget domain.Budget) error {
	if err := budget.Validate(); err != nil {
		return err
	}

	return s.budgetRepo.Upsert(ctx, budget)
}

// GetBudgets возвращает все бюджеты
func (s *LedgerService) GetBudgets(ctx context.Context) ([]domain.Budget, error) {
	return s.budgetRepo.List(ctx)
}

// AddTransaction добавляет новую транзакцию с проверкой бюджета
func (s *LedgerService) AddTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	if err := tx.Validate(); err != nil {
		return tx, err
	}

	// Проверка бюджета
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

	// Создание транзакции
	id, err := s.transactionRepo.Create(ctx, tx)
	if err != nil {
		return tx, err
	}

	tx.ID = id
	return tx, nil
}

// GetTransactions возвращает все транзакции
func (s *LedgerService) GetTransactions(ctx context.Context) ([]domain.Transaction, error) {
	return s.transactionRepo.List(ctx)
}

// GetReportSummary возвращает сводку расходов за период
func (s *LedgerService) GetReportSummary(ctx context.Context, from, to time.Time) ([]domain.ReportSummary, error) {
	return s.transactionRepo.GetReportSummary(ctx, from, to)
}

