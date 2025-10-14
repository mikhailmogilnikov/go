package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

var (
	transactions []Transaction
	budgets      map[string]Budget
)

func init() {
	budgets = make(map[string]Budget)
}

func SetBudget(b Budget) {
	budgets[b.Category] = b
}

func getCategoryTotal(category string) float64 {
	var total float64
	for _, tx := range transactions {
		if tx.Category == category {
			total += tx.Amount
		}
	}
	return total
}

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("сумма транзакции не может быть равна 0")
	}

	if budget, exists := budgets[tx.Category]; exists {
		currentTotal := getCategoryTotal(tx.Category)
		newTotal := currentTotal + tx.Amount

		if newTotal > budget.Limit {
			return errors.New("budget exceeded")
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func LoadBudgets(r io.Reader) error {
	decoder := json.NewDecoder(r)

	var budgetList []Budget
	if err := decoder.Decode(&budgetList); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	for _, budget := range budgetList {
		if budget.Category == "" {
			return errors.New("категория бюджета не может быть пустой")
		}
		if budget.Limit <= 0 {
			return fmt.Errorf("лимит бюджета для категории '%s' должен быть положительным числом", budget.Category)
		}
		SetBudget(budget)
	}

	return nil
}

func main() {
	fmt.Println("Ledger service started")
	fmt.Println()

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Printf("Не удалось открыть файл budgets.json: %v\n", err)
	} else {
		reader := bufio.NewReader(file)
		if err := LoadBudgets(reader); err != nil {
			fmt.Printf("Ошибка при загрузке бюджетов: %v\n", err)
		} else {
			fmt.Println("Бюджеты загружены из файла")
		}
		file.Close()
	}

	SetBudget(Budget{
		Category: "Развлечения",
		Limit:    3000.00,
		Period:   "месяц",
	})
	SetBudget(Budget{
		Category: "Транспорт",
		Limit:    10000.00,
		Period:   "месяц",
	})

	fmt.Println()
	fmt.Println("Добавление транзакций:")

	tx1 := Transaction{
		Amount:      2000.00,
		Category:    "Еда",
		Description: "Покупка продуктов",
		Date:        time.Now(),
	}

	tx2 := Transaction{
		Amount:      1500.00,
		Category:    "Еда",
		Description: "Ресторан",
		Date:        time.Now(),
	}

	tx3 := Transaction{
		Amount:      2000.00,
		Category:    "Еда",
		Description: "Превышение бюджета",
		Date:        time.Now(),
	}

	tx4 := Transaction{
		Amount:      800.00,
		Category:    "Развлечения",
		Description: "Кино",
		Date:        time.Now(),
	}

	tx5 := Transaction{
		Amount:      0,
		Category:    "Тест",
		Description: "Нулевая сумма",
		Date:        time.Now(),
	}

	if err := AddTransaction(tx1); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции: %v\n", err)
	} else {
		fmt.Println("Транзакция 1 добавлена")
	}

	if err := AddTransaction(tx2); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции: %v\n", err)
	} else {
		fmt.Println("Транзакция 2 добавлена")
	}

	if err := AddTransaction(tx3); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции: %v\n", err)
	} else {
		fmt.Println("Транзакция 3 добавлена")
	}

	if err := AddTransaction(tx4); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции: %v\n", err)
	} else {
		fmt.Println("Транзакция 4 добавлена")
	}

	if err := AddTransaction(tx5); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции: %v\n", err)
	} else {
		fmt.Println("Транзакция 5 добавлена")
	}

	fmt.Println()
	fmt.Println("Список транзакций:")
	allTransactions := ListTransactions()
	for _, tx := range allTransactions {
		fmt.Printf("ID: %d, Сумма: %.2f, Категория: %s\n", tx.ID, tx.Amount, tx.Category)
	}

	fmt.Println()
	fmt.Println("Статус бюджетов:")
	for category, budget := range budgets {
		spent := getCategoryTotal(category)
		fmt.Printf("%s: потрачено %.2f из %.2f\n", category, spent, budget.Limit)
	}
}
