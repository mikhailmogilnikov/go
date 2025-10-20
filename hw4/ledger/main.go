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

type Validatable interface {
	Validate() error
}

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

func (t Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("сумма транзакции должна быть положительным числом")
	}
	if t.Category == "" {
		return errors.New("категория транзакции не может быть пустой")
	}
	return nil
}

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("лимит бюджета должен быть положительным числом")
	}
	if b.Category == "" {
		return errors.New("категория бюджета не может быть пустой")
	}
	return nil
}

var (
	transactions []Transaction
	budgets      map[string]Budget
)

func init() {
	budgets = make(map[string]Budget)
}

func CheckValid(v Validatable) error {
	return v.Validate()
}

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации бюджета: %w", err)
	}
	budgets[b.Category] = b
	return nil
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
	if err := tx.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации транзакции: %w", err)
	}

	if budget, exists := budgets[tx.Category]; exists {
		currentTotal := getCategoryTotal(tx.Category)
		newTotal := currentTotal + tx.Amount

		if newTotal > budget.Limit {
			return errors.New("превышен лимит бюджета")
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
		if err := SetBudget(budget); err != nil {
			return err
		}
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

	budget1 := Budget{
		Category: "Развлечения",
		Limit:    3000.00,
		Period:   "месяц",
	}
	budget2 := Budget{
		Category: "Транспорт",
		Limit:    10000.00,
		Period:   "месяц",
	}

	if err := SetBudget(budget1); err != nil {
		fmt.Printf("Ошибка при добавлении бюджета: %v\n", err)
	}
	if err := SetBudget(budget2); err != nil {
		fmt.Printf("Ошибка при добавлении бюджета: %v\n", err)
	}

	fmt.Println()
	fmt.Println("=== Демонстрация полиморфизма через интерфейс Validatable ===")

	testTx1 := Transaction{
		Amount:      1000.00,
		Category:    "Еда",
		Description: "Тест валидации",
		Date:        time.Now(),
	}
	fmt.Printf("Валидация корректной транзакции: ")
	if err := CheckValid(testTx1); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	testTx2 := Transaction{
		Amount:      1000.00,
		Category:    "",
		Description: "Тест валидации",
		Date:        time.Now(),
	}
	fmt.Printf("Валидация транзакции с пустой категорией: ")
	if err := CheckValid(testTx2); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	testTx3 := Transaction{
		Amount:      -500.00,
		Category:    "Еда",
		Description: "Тест валидации",
		Date:        time.Now(),
	}
	fmt.Printf("Валидация транзакции с отрицательной суммой: ")
	if err := CheckValid(testTx3); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	testBudget1 := Budget{
		Category: "Тест",
		Limit:    5000.00,
		Period:   "месяц",
	}
	fmt.Printf("Валидация корректного бюджета: ")
	if err := CheckValid(testBudget1); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	testBudget2 := Budget{
		Category: "",
		Limit:    5000.00,
		Period:   "месяц",
	}
	fmt.Printf("Валидация бюджета с пустой категорией: ")
	if err := CheckValid(testBudget2); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	testBudget3 := Budget{
		Category: "Тест",
		Limit:    -1000.00,
		Period:   "месяц",
	}
	fmt.Printf("Валидация бюджета с отрицательным лимитом: ")
	if err := CheckValid(testBudget3); err != nil {
		fmt.Printf("ОШИБКА - %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Println()
	fmt.Println("=== Добавление транзакций ===")

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

	tx6 := Transaction{
		Amount:      -100.00,
		Category:    "Тест",
		Description: "Отрицательная сумма",
		Date:        time.Now(),
	}

	tx7 := Transaction{
		Amount:      500.00,
		Category:    "",
		Description: "Пустая категория",
		Date:        time.Now(),
	}

	if err := AddTransaction(tx1); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 1: %v\n", err)
	} else {
		fmt.Println("Транзакция 1 добавлена")
	}

	if err := AddTransaction(tx2); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 2: %v\n", err)
	} else {
		fmt.Println("Транзакция 2 добавлена")
	}

	if err := AddTransaction(tx3); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 3: %v\n", err)
	} else {
		fmt.Println("Транзакция 3 добавлена")
	}

	if err := AddTransaction(tx4); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 4: %v\n", err)
	} else {
		fmt.Println("Транзакция 4 добавлена")
	}

	if err := AddTransaction(tx5); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 5: %v\n", err)
	} else {
		fmt.Println("Транзакция 5 добавлена")
	}

	if err := AddTransaction(tx6); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 6: %v\n", err)
	} else {
		fmt.Println("Транзакция 6 добавлена")
	}

	if err := AddTransaction(tx7); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 7: %v\n", err)
	} else {
		fmt.Println("Транзакция 7 добавлена")
	}

	fmt.Println()
	fmt.Println("=== Список транзакций ===")
	allTransactions := ListTransactions()
	for _, tx := range allTransactions {
		fmt.Printf("ID: %d, Сумма: %.2f, Категория: %s, Описание: %s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description)
	}

	fmt.Println()
	fmt.Println("=== Статус бюджетов ===")
	for category, budget := range budgets {
		spent := getCategoryTotal(category)
		remaining := budget.Limit - spent
		fmt.Printf("%s: потрачено %.2f из %.2f (осталось %.2f)\n",
			category, spent, budget.Limit, remaining)
	}
}
