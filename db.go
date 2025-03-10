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

// Создаёт таблицы в базе данных
func CreateTable(db *sql.DB) error {
	// Создание таблицы для финансовых операций
	queryTransactions := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operation TEXT,
		type TEXT,
		account TEXT,
		amount REAL,
		before_operation REAL,
		after_operation REAL
	)`

	_, err := db.Exec(queryTransactions)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса для создания таблицы transactions: %w", err)
	}

	// Создание таблицы для настроек (например, цвет кружочка)
	querySettings := `CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,   -- Ключ настройки (например, "circle_color")
		value TEXT NOT NULL     -- Значение настройки (например, "#a19f9f")
	)`

	_, err = db.Exec(querySettings)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса для создания таблицы settings: %w", err)
	}

	return nil
}

// Сохранить настройку в базу данных
func SaveSetting(db *sql.DB, key, value string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	return err
}

// Загрузить настройку из базы данных
func LoadSetting(db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// Обновляет данные в базе данных
func UpdateDatabase(db *sql.DB, operation, typeOp, account string, amount, before, after float64) error {
	_, err := db.Exec(`INSERT INTO transactions (operation, type, account, amount, before_operation, after_operation) VALUES (?, ?, ?, ?, ?, ?)`, operation, typeOp, account, amount, before, after)
	return err
}
