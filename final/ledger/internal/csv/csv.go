package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
)

func ParseCSV(data []byte, userID int64) ([]domain.Transaction, []string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	var transactions []domain.Transaction
	var errors []string

	for i, record := range records {
		if i == 0 && (record[0] == "amount" || record[0] == "Amount" || record[0] == "сумма" || record[0] == "Сумма") {
			continue
		}

		if len(record) < 2 {
			errors = append(errors, fmt.Sprintf("row %d: not enough columns", i+1))
			continue
		}

		amount, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: invalid amount '%s'", i+1, record[0]))
			continue
		}

		category := record[1]
		if category == "" {
			errors = append(errors, fmt.Sprintf("row %d: empty category", i+1))
			continue
		}

		description := ""
		if len(record) > 2 {
			description = record[2]
		}

		date := time.Now()
		if len(record) > 3 && record[3] != "" {
			parsed, err := time.Parse("2006-01-02", record[3])
			if err != nil {
				errors = append(errors, fmt.Sprintf("row %d: invalid date '%s', using today", i+1, record[3]))
			} else {
				date = parsed
			}
		}

		tx := domain.Transaction{
			UserID:      userID,
			Amount:      amount,
			Category:    category,
			Description: description,
			Date:        date,
		}

		if err := tx.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}

		transactions = append(transactions, tx)
	}

	return transactions, errors, nil
}

func GenerateCSV(transactions []domain.Transaction) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"amount", "category", "description", "date"}); err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		record := []string{
			fmt.Sprintf("%.2f", tx.Amount),
			tx.Category,
			tx.Description,
			tx.Date.Format("2006-01-02"),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}



