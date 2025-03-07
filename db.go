package main

import (
	"database/sql"
	"fmt"
)

// Открывает соединение с базой данных
func OpenDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite", "./base.db")
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}
	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка соединения с базой данных: %w", err)
	}
	return database, nil
}

// Создаёт таблицу в базе данных
func CreateTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY,
		operation TEXT,
		amount FLOAT64,
		operation_type TEXT,
		count_type TEXT,
		before_operation FLOAT64,
		after_operation FLOAT64
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса для создания таблицы: %w", err)
	}
	return nil
}

// Обновляет данные в базе данных
func UpdateDatabase(db *sql.DB, operation, operationType, countType string, amount float64, beforeOperation, afterOperation float64) error {
	// Вставка новой записи
	query := `INSERT INTO transactions (operation, amount, operation_type, count_type, before_operation, after_operation) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, operation, amount, operationType, countType, beforeOperation, afterOperation)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса для обновления базы данных: %w", err)
	}
	return nil
}
