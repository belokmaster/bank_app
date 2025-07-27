import sqlite3
from datetime import datetime


DB_NAME = 'finance_tracker.db'

INCOME_CATEGORIES = ["Зарплата", "Подработка", "Подарок", "Проценты", "Другое"]
EXPENSE_CATEGORIES = ["Еда", "Транспорт", "Жилье", "Развлечения", "Здоровье", "Покупки", "Счета", "Другое"]

def init_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute('''
        CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT NOT NULL,
            category TEXT NOT NULL,
            amount REAL NOT NULL,
            description TEXT,
            date TEXT NOT NULL
        )
    ''')

    conn.commit()
    conn.close()


def add_transaction_db(trans_type, category, amount, description):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    current_date = datetime.now().strftime('%Y-%m-%d')
    
    cursor.execute('INSERT INTO transactions (type, category, amount, description, date) VALUES (?, ?, ?, ?, ?)',
                   (trans_type, category, amount, description, current_date))
    conn.commit()
    conn.close()


def get_all_transactions_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute('SELECT id, type, category, amount, description, date FROM transactions ORDER BY date DESC, id DESC')
    transactions = cursor.fetchall()
    conn.close()

    return transactions


def delete_transaction_db(transaction_id):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute('DELETE FROM transactions WHERE id = ?', (transaction_id,))
    
    conn.commit()
    conn.close()


def calculate_totals_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute("SELECT SUM(amount) FROM transactions WHERE type = 'Доход'")
    total_income = cursor.fetchone()[0] or 0
    cursor.execute("SELECT SUM(amount) FROM transactions WHERE type = 'Расход'")
    total_expense = cursor.fetchone()[0] or 0

    conn.close()
    balance = total_income - total_expense
    return {'income': total_income, 'expense': total_expense, 'balance': balance}