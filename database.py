import sqlite3
from datetime import datetime


DB_NAME = 'finance_tracker.db'


def init_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute("PRAGMA foreign_keys = ON;")

    cursor.execute('''
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            type TEXT NOT NULL  -- 'Доход' или 'Расход'
        )
    ''')

    cursor.execute('''
        CREATE TABLE IF NOT EXISTS subcategories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            category_id INTEGER NOT NULL,
            FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
        )
    ''')


    cursor.execute('DROP TABLE IF EXISTS transactions')
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT NOT NULL,
            amount REAL NOT NULL,
            description TEXT,
            date TEXT NOT NULL,
            category_id INTEGER,
            subcategory_id INTEGER,
            FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL,
            FOREIGN KEY (subcategory_id) REFERENCES subcategories (id) ON DELETE SET NULL
        )
    ''')

    conn.commit()
    conn.close()


def add_category_db(name, type):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    try:
        cursor.execute('INSERT INTO categories (name, type) VALUES (?, ?)', (name, type))
        conn.commit()
    except sqlite3.IntegrityError:
        print(f"Категория '{name}' уже существует.")
    finally:
        conn.close()

def add_subcategory_db(name, category_id):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute('INSERT INTO subcategories (name, category_id) VALUES (?, ?)', (name, category_id))

    conn.commit()
    conn.close()

def get_categories_with_subcategories_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()
    categories_dict = {'Доход': [], 'Расход': []}
    
    cursor.execute('SELECT id, name, type FROM categories ORDER BY name')
    categories = cursor.fetchall()
    
    for cat_id, cat_name, cat_type in categories:
        cursor.execute('SELECT id, name FROM subcategories WHERE category_id = ? ORDER BY name', (cat_id,))
        subcategories = cursor.fetchall()
        
        category_info = {
            'id': cat_id,
            'name': cat_name,
            'subcategories': [{'id': sub_id, 'name': sub_name} for sub_id, sub_name in subcategories]
        }
        if cat_type in categories_dict:
            categories_dict[cat_type].append(category_info)
            
    conn.close()
    return categories_dict

def add_transaction_db(trans_type, category_id, subcategory_id, amount, description):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    current_date = datetime.now().strftime('%Y-%m-%d')
    cursor.execute('INSERT INTO transactions (type, amount, description, date, category_id, subcategory_id) VALUES (?, ?, ?, ?, ?, ?)',
                   (trans_type, amount, description, current_date, category_id, subcategory_id))
    conn.commit()
    conn.close()

def get_all_transactions_db():
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    query = '''
        SELECT 
            t.id, t.type, c.name, s.name, t.amount, t.description, t.date
        FROM transactions t
        LEFT JOIN categories c ON t.category_id = c.id
        LEFT JOIN subcategories s ON t.subcategory_id = s.id
        ORDER BY t.date DESC, t.id DESC
    '''

    cursor.execute(query)
    transactions = cursor.fetchall()
    conn.close()

    return transactions

def delete_transaction_db(transaction_id):
    conn = sqlite3.connect(DB_NAME)
    cursor = conn.cursor()

    cursor.execute("PRAGMA foreign_keys = ON;")
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