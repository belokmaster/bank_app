package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

func getCurrentBalance(db *sql.DB) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT after_operation FROM transactions ORDER BY id DESC LIMIT 1").Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return balance, err
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		balance, _ := getCurrentBalance(db)
		tmpl.Execute(w, balance)
	})

	http.HandleFunc("/deposit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			http.Error(w, "Некорректная сумма", http.StatusBadRequest)
			return
		}
		typeOp := r.FormValue("type")
		countType := r.FormValue("account")
		before, err := getCurrentBalance(db)
		if err != nil {
			http.Error(w, "Ошибка получения баланса", http.StatusInternalServerError)
			return
		}

		var after float64
		if typeOp == "income" {
			after = before + amount
		} else if typeOp == "expense" {
			after = before - amount
		} else {
			http.Error(w, "Некорректный тип операции", http.StatusBadRequest)
			return
		}

		UpdateDatabase(db, "+", typeOp, countType, amount, before, after)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
