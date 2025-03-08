package main

import (
	"database/sql"
	"fmt"
)

// Открывает соединение с базой данных
func OpenDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite", "./balance.db")
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operation TEXT,
		type TEXT,
		account TEXT,
		amount REAL,
		before_operation REAL,
		after_operation REAL
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса для создания таблицы: %w", err)
	}
	return nil
}

// Обновляет данные в базе данных
func UpdateDatabase(db *sql.DB, operation, typeOp, account string, amount, before, after float64) error {
	_, err := db.Exec(`INSERT INTO transactions (operation, type, account, amount, before_operation, after_operation) VALUES (?, ?, ?, ?, ?, ?)`, operation, typeOp, account, amount, before, after)
	return err
}
