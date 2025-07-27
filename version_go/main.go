package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

func createProfitLossPlot(history []map[string]interface{}) (*plot.Plot, error) {
	p := plot.New()

	p.Title.Text = "Прибыль и расходы"
	p.X.Label.Text = "Операции"
	p.Y.Label.Text = "Сумма"

	profitPoints := make(plotter.XYs, 0)
	lossPoints := make(plotter.XYs, 0)

	for i, entry := range history {
		if entry["type"].(string) == "income" {
			profitPoints = append(profitPoints, plotter.XY{X: float64(i), Y: entry["amount"].(float64)})
		} else if entry["type"].(string) == "expense" {
			lossPoints = append(lossPoints, plotter.XY{X: float64(i), Y: entry["amount"].(float64)})
		}
	}

	profitLine, err := plotter.NewLine(profitPoints)
	if err != nil {
		return nil, err
	}
	profitLine.Color = plotutil.Color(0)

	lossLine, err := plotter.NewLine(lossPoints)
	if err != nil {
		return nil, err
	}
	lossLine.Color = plotutil.Color(1)

	p.Add(profitLine, lossLine)
	p.Legend.Add("Прибыль", profitLine)
	p.Legend.Add("Расходы", lossLine)

	return p, nil
}

func main() {
	// Открытие базы данных
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблиц
	err = CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	// Загрузка цвета кружочков из базы данных
	circleColor, err := LoadSetting(db, "circle_color")
	if err != nil {
		log.Fatal("Ошибка загрузки цвета кружочков:", err)
	}
	if circleColor == "" {
		circleColor = "#a19f9f" // Значение по умолчанию
		log.Println("Используется цвет по умолчанию:", circleColor)
	} else {
		log.Println("Загружен цвет из базы данных:", circleColor)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		balance, _ := getCurrentBalance(db)
		data := struct {
			Balance float64
			Color   string
		}{
			Balance: balance,
			Color:   circleColor,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		color := r.FormValue("color")
		if color != "" {
			// Сохраняем новый цвет в базу данных
			err := SaveSetting(db, "circle_color", color)
			if err != nil {
				log.Printf("Ошибка сохранения цвета: %v", err)
				http.Error(w, "Ошибка сохранения цвета", http.StatusInternalServerError)
				return
			}
			log.Printf("Цвет успешно сохранен в базу данных: %s", color)

			// Обновляем цвет в памяти
			circleColor = color
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/category", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		category := r.FormValue("category")
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

		// Сохраняем транзакцию с категорией
		UpdateDatabaseWithCategory(db, typeOp, countType, amount, before, after, category)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		category := r.FormValue("category") // Получаем категорию из формы
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

		UpdateDatabase(db, category, typeOp, countType, amount, before, after)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		history, err := GetTransactionHistory(db)
		if err != nil {
			log.Printf("Ошибка получения истории: %v", err)
			http.Error(w, "Ошибка получения истории операций", http.StatusInternalServerError)
			return
		}

		// Логируем данные перед отправкой
		log.Printf("Отправляемые данные: %+v", history)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(history); err != nil {
			log.Printf("Ошибка кодирования JSON: %v", err)
			http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/graphs", func(w http.ResponseWriter, r *http.Request) {
		// Получаем историю операций из базы данных
		history, err := GetTransactionHistory(db)
		if err != nil {
			log.Printf("Ошибка получения истории: %v", err)
			http.Error(w, "Ошибка получения истории операций", http.StatusInternalServerError)
			return
		}

		// Логируем данные перед отправкой
		log.Printf("Отправляемые данные: %+v", history)

		// Создаем график на основе истории
		p, err := createProfitLossPlot(history)
		if err != nil {
			log.Printf("Ошибка создания графика: %v", err)
			http.Error(w, "Ошибка создания графика", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовок для ответа в формате PNG
		w.Header().Set("Content-Type", "image/png")

		// Создаем PNG-изображение и отправляем его в ответ
		writer, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
		if err != nil {
			log.Printf("Ошибка создания изображения: %v", err)
			http.Error(w, "Ошибка создания изображения", http.StatusInternalServerError)
			return
		}

		_, err = writer.WriteTo(w)
		if err != nil {
			log.Printf("Ошибка записи изображения: %v", err)
			http.Error(w, "Ошибка записи изображения", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
