package domain

import (
	"testing"
	"time"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				UserID:   1,
				Amount:   100.50,
				Category: "food",
				Date:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			tx: Transaction{
				UserID:   1,
				Amount:   0,
				Category: "food",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			tx: Transaction{
				UserID:   1,
				Amount:   -50,
				Category: "food",
			},
			wantErr: true,
		},
		{
			name: "empty category",
			tx: Transaction{
				UserID: 1,
				Amount: 100,
			},
			wantErr: true,
		},
		{
			name: "zero user_id",
			tx: Transaction{
				Amount:   100,
				Category: "food",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}



