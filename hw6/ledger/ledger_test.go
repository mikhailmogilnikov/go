package ledger

import (
	"testing"
	"time"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				Amount:   100.0,
				Category: "еда",
				Date:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			tx: Transaction{
				Amount:   0,
				Category: "еда",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "сумма транзакции должна быть положительным числом",
		},
		{
			name: "negative amount",
			tx: Transaction{
				Amount:   -50.0,
				Category: "еда",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "сумма транзакции должна быть положительным числом",
		},
		{
			name: "empty category",
			tx: Transaction{
				Amount:   100.0,
				Category: "",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "категория транзакции не может быть пустой",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Transaction.Validate() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Transaction.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Transaction.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBudget_Validate(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid budget",
			budget: Budget{
				Category: "еда",
				Limit:    5000.0,
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: Budget{
				Category: "еда",
				Limit:    0,
			},
			wantErr: true,
			errMsg:  "лимит бюджета должен быть положительным числом",
		},
		{
			name: "negative limit",
			budget: Budget{
				Category: "еда",
				Limit:    -1000.0,
			},
			wantErr: true,
			errMsg:  "лимит бюджета должен быть положительным числом",
		},
		{
			name: "empty category",
			budget: Budget{
				Category: "",
				Limit:    5000.0,
			},
			wantErr: true,
			errMsg:  "категория бюджета не может быть пустой",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.budget.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Budget.Validate() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Budget.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Budget.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBudgetAndTransaction(t *testing.T) {
	t.Run("within budget", func(t *testing.T) {
		Reset()
		defer Reset()

		budget := Budget{Category: "еда", Limit: 5000.0}
		if err := SetBudget(budget); err != nil {
			t.Fatalf("SetBudget() error = %v", err)
		}

		tx := Transaction{
			Amount:   3000.0,
			Category: "еда",
			Date:     time.Now(),
		}

		if err := AddTransaction(tx); err != nil {
			t.Errorf("AddTransaction() error = %v, want nil", err)
		}

		transactions := ListTransactions()
		if len(transactions) != 1 {
			t.Errorf("ListTransactions() len = %d, want 1", len(transactions))
		}
	})

	t.Run("budget exceeded", func(t *testing.T) {
		Reset()
		defer Reset()

		budget := Budget{Category: "еда", Limit: 5000.0}
		if err := SetBudget(budget); err != nil {
			t.Fatalf("SetBudget() error = %v", err)
		}

		tx1 := Transaction{
			Amount:   3000.0,
			Category: "еда",
			Date:     time.Now(),
		}

		if err := AddTransaction(tx1); err != nil {
			t.Fatalf("AddTransaction() error = %v", err)
		}

		tx2 := Transaction{
			Amount:   2500.0,
			Category: "еда",
			Date:     time.Now(),
		}

		if err := AddTransaction(tx2); err == nil {
			t.Error("AddTransaction() expected error, got nil")
		} else if err.Error() != "budget exceeded" {
			t.Errorf("AddTransaction() error = %v, want 'budget exceeded'", err)
		}

		transactions := ListTransactions()
		if len(transactions) != 1 {
			t.Errorf("ListTransactions() len = %d, want 1", len(transactions))
		}
	})
}

