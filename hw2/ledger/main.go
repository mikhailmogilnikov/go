package main

import (
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("сумма транзакции не может быть равна 0")
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

func main() {
	fmt.Println("Ledger service started")
	fmt.Println()

	// Добавление тестовых транзакций
	tx1 := Transaction{
		Amount:      1500.50,
		Category:    "Продукты",
		Description: "Покупка продуктов в супермаркете",
		Date:        time.Now(),
	}

	tx2 := Transaction{
		Amount:      3200.00,
		Category:    "Транспорт",
		Description: "Заправка автомобиля",
		Date:        time.Now().Add(-24 * time.Hour),
	}

	tx3 := Transaction{
		Amount:      850.75,
		Category:    "Развлечения",
		Description: "Билеты в кино",
		Date:        time.Now().Add(-2 * 24 * time.Hour),
	}

	// Пробуем добавить транзакцию с нулевой суммой
	tx4 := Transaction{
		Amount:      0,
		Category:    "Тест",
		Description: "Тестовая транзакция с нулевой суммой",
		Date:        time.Now(),
	}

	fmt.Println("Добавление транзакций:")
	if err := AddTransaction(tx1); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 1: %v\n", err)
	} else {
		fmt.Println("✓ Транзакция 1 успешно добавлена")
	}

	if err := AddTransaction(tx2); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 2: %v\n", err)
	} else {
		fmt.Println("✓ Транзакция 2 успешно добавлена")
	}

	if err := AddTransaction(tx3); err != nil {
		fmt.Printf("Ошибка при добавлении транзакции 3: %v\n", err)
	} else {
		fmt.Println("✓ Транзакция 3 успешно добавлена")
	}

	// Попытка добавить транзакцию с нулевой суммой
	if err := AddTransaction(tx4); err != nil {
		fmt.Printf("✗ Ошибка при добавлении транзакции 4: %v\n", err)
	} else {
		fmt.Println("✓ Транзакция 4 успешно добавлена")
	}

	fmt.Println("\n" + "Список всех транзакций:")
	fmt.Println("--------------------------------------------------")

	allTransactions := ListTransactions()
	if len(allTransactions) == 0 {
		fmt.Println("Нет транзакций")
	} else {
		for _, tx := range allTransactions {
			fmt.Printf("ID: %d\n", tx.ID)
			fmt.Printf("Сумма: %.2f руб.\n", tx.Amount)
			fmt.Printf("Категория: %s\n", tx.Category)
			fmt.Printf("Описание: %s\n", tx.Description)
			fmt.Printf("Дата: %s\n", tx.Date.Format("02.01.2006 15:04:05"))
			fmt.Println("--------------------------------------------------")
		}
		fmt.Printf("\nВсего транзакций: %d\n", len(allTransactions))

		// Подсчёт общей суммы
		var total float64
		for _, tx := range allTransactions {
			total += tx.Amount
		}
		fmt.Printf("Общая сумма: %.2f руб.\n", total)
	}
}
