package domain

import "testing"

func TestBudget_Validate(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name: "valid monthly budget",
			budget: Budget{
				UserID:      1,
				Category:    "food",
				LimitAmount: 15000,
				Period:      "monthly",
			},
			wantErr: false,
		},
		{
			name: "valid weekly budget",
			budget: Budget{
				UserID:      1,
				Category:    "transport",
				LimitAmount: 5000,
				Period:      "weekly",
			},
			wantErr: false,
		},
		{
			name: "default period",
			budget: Budget{
				UserID:      1,
				Category:    "food",
				LimitAmount: 10000,
				Period:      "",
			},
			wantErr: false, // period becomes "monthly"
		},
		{
			name: "zero limit",
			budget: Budget{
				UserID:   1,
				Category: "food",
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			budget: Budget{
				UserID:      1,
				Category:    "food",
				LimitAmount: -100,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			budget: Budget{
				UserID:      1,
				LimitAmount: 10000,
			},
			wantErr: true,
		},
		{
			name: "zero user_id",
			budget: Budget{
				Category:    "food",
				LimitAmount: 10000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.budget.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}



