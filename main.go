package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// Функция для получения текущего баланса из базы данных
func GetCurrentBalance(db *sql.DB) (float64, error) {
	var balance float64
	query := "SELECT after_operation FROM transactions ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(query).Scan(&balance)
	if err != nil {
		// Если ошибка, то возвращаем 0 как начальный баланс
		if err.Error() == "no rows in result set" {
			return 0, nil // Таблица пуста, начинаем с 0
		}
		return 0, err
	}
	return balance, nil
}

func main() {
	// Открытие базы данных
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы
	err = CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	// Чтение пользовательского ввода
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("1 - внести средства на счет, 2 - выйти")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var operation string
		switch input {
		case "1":
			operation = "+"

			fmt.Println("Введите сумму:")
			amountInput, _ := reader.ReadString('\n')
			amountInput = strings.TrimSpace(amountInput)

			amount, err := strconv.ParseFloat(amountInput, 64)
			if err != nil {
				fmt.Println("Ошибка: введите число для суммы")
				continue
			}

			fmt.Println("Введите тип операции:")
			operationType, _ := reader.ReadString('\n')
			operationType = strings.TrimSpace(operationType)

			fmt.Println("Введите счет (например, 'Наличные' или 'Банковский перевод'):")
			countType, _ := reader.ReadString('\n')
			countType = strings.TrimSpace(countType)

			// Получаем текущее значение баланса
			beforeOperation, err := GetCurrentBalance(db)
			if err != nil {
				fmt.Println("Ошибка при получении текущего баланса:", err)
				continue
			}

			// Если таблица пуста, то начнем с нуля
			if beforeOperation == 0 {
				beforeOperation = 0.0
			}

			// Вычисляем новый баланс после операции
			afterOperation := beforeOperation + amount

			// Обновление базы данных
			err = UpdateDatabase(db, operation, operationType, countType, amount, beforeOperation, afterOperation)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Запись успешно добавлена!")

		case "2":
			// Выход из программы
			fmt.Println("Выход из программы.")
			return

		default:
			fmt.Println("Неизвестная команда, попробуйте снова.")
		}
	}
}
