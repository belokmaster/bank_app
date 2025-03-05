package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := "money.xlsx"

	f, sheetName, err := openExcelFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeExcelFile(f)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("1 - добавить запись, 2 - снять деньги со счета, q - выход")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Введите сумму для зачисления средств.")
			sumInput, _ := reader.ReadString('\n')
			sumInput = strings.TrimSpace(sumInput)

			if _, err := strconv.ParseFloat(sumInput, 64); err != nil {
				fmt.Println("Ошибка: введите только число (например, 1 или 1.05)")
			}

			fmt.Println("Введите тип операции.")
			operationInput, _ := reader.ReadString('\n')

			err := updateExcelCell(f, sheetName, sumInput, "+", operationInput)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = saveExcelFile(f, fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Данные успешно обновлены!")
			}
		case "2":
			fmt.Println("Введите сумму для снятия средств.")
			sumInput, _ := reader.ReadString('\n')
			sumInput = strings.TrimSpace(sumInput)

			if _, err := strconv.ParseFloat(sumInput, 64); err != nil {
				fmt.Println("Ошибка: введите только число (например, 1 или 1.05)")
			}

			fmt.Println("Введите тип операции.")
			operationInput, _ := reader.ReadString('\n')

			err := updateExcelCell(f, sheetName, sumInput, "-", operationInput)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = saveExcelFile(f, fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Данные успешно обновлены!")
			}

		case "q":
			fmt.Println("Работа программы завершена.")
			return
		default:
			fmt.Println("Неизвестная команда, попробуйте снова.")
		}
	}
}
