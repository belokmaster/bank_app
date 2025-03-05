package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Открывает Excel-файл и возвращает объект файла и имя первого листа
func openExcelFile(fileName string) (*excelize.File, string, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, "", fmt.Errorf("не удалось открыть файл: %w", err)
	}

	// Получаем имя первого листа
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, "", fmt.Errorf("файл не содержит листов")
	}
	return f, sheets[0], nil
}

// Получает номер последней заполненной строки в столбце B
func getLastRow(f *excelize.File, sheet, column string) (int, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	lastRow := 1 // Если столбец пуст, начнем с первой строки
	for i, row := range rows {
		if len(row) >= 2 && row[1] != "" { // Проверяем столбец B (индекс 1)
			lastRow = i + 1 // Excel начинается с 1, а массив с 0
		}
	}
	return lastRow + 1, nil // Следующая строка после последней заполненной
}

func updateExcelCell(f *excelize.File, sheet, value, operation, typeOperation string) error {
	lastRow, err := getLastRow(f, sheet, "B")
	if err != nil {
		return err
	}

	// Формируем адреса ячеек
	date := time.Now().Format("2006-01-02")
	cellA := fmt.Sprintf("A%d", lastRow)
	cellB := fmt.Sprintf("B%d", lastRow)
	cellC := fmt.Sprintf("C%d", lastRow)
	cellD := fmt.Sprintf("D%d", lastRow)
	cellE := fmt.Sprintf("E%d", lastRow)
	cellF := fmt.Sprintf("F%d", lastRow)

	// Формируем адреса ячеек из предыдущей строки
	prevRow := lastRow - 1
	cellFPrev := fmt.Sprintf("F%d", prevRow)

	// Записываем дату и значение
	if err := f.SetCellValue(sheet, cellA, date); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellA, err)
	}
	if err := f.SetCellValue(sheet, cellB, value); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellB, err)
	}
	if err := f.SetCellValue(sheet, cellC, operation); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellC, err)
	}
	if err := f.SetCellValue(sheet, cellD, typeOperation); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellD, err)
	}

	// Читаем значение из F предыдущей строки
	valFPrev, err := f.GetCellValue(sheet, cellFPrev)
	if err != nil || valFPrev == "" {
		valFPrev = "0" // Если пустая ячейка или ошибка чтения, считаем 0
	}

	// Записываем значение F прошлого столбца в E текущего столбца
	if err := f.SetCellValue(sheet, cellE, valFPrev); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellE, err)
	}

	// Преобразуем valFPrev и value в числа и складываем
	numE, _ := strconv.ParseFloat(valFPrev, 64)
	value = strings.TrimSpace(value)
	numB, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования строки %s: %v\n", value, err)
		numB = 0 // Если ошибка, устанавливаем 0
	}

	var sum float64
	if operation == "+" {
		sum = numE + numB
	} else {
		sum = numE - numB
	}

	// Записываем сумму в F
	if err := f.SetCellValue(sheet, cellF, sum); err != nil {
		return fmt.Errorf("не удалось записать в ячейку %s: %w", cellF, err)
	}

	return nil
}

// Закрытие файла Excel
func closeExcelFile(f *excelize.File) {
	if err := f.Close(); err != nil {
		fmt.Println("Ошибка при закрытии файла:", err)
	}
}

// Сохранение Excel-файла
func saveExcelFile(f *excelize.File, fileName string) error {
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("ошибка сохранения файла: %w", err)
	}
	return nil
}
