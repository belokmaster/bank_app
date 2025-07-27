from flask import Flask, render_template, request, redirect, url_for
import database as db

app = Flask(__name__)

# маршрутизаторы
@app.route('/')
def index():
    transactions = db.get_all_transactions_db()
    totals = db.calculate_totals_db()
    return render_template('index.html', 
                           transactions=transactions, 
                           totals=totals,
                           income_categories=db.INCOME_CATEGORIES,
                           expense_categories=db.EXPENSE_CATEGORIES)


@app.route('/add', methods=['POST'])
def add():
    trans_type = request.form['type']
    category = request.form['category']
    description = request.form['description']
    
    try:
        amount = float(request.form['amount'])
        if amount <= 0:
            return "Ошибка: сумма должна быть положительным числом. <a href='/'>Назад</a>"
    except ValueError:
        return "Ошибка: сумма должна быть числом. <a href='/'>Назад</a>"

    db.add_transaction_db(trans_type, category, amount, description)
    return redirect(url_for('index'))


@app.route('/delete/<int:transaction_id>', methods=['POST'])
def delete(transaction_id):
    db.delete_transaction_db(transaction_id)
    return redirect(url_for('index'))

if __name__ == '__main__':
    db.init_db()
    try:
        app.run(host='0.0.0.0', port=80, debug=True)
    except OSError as e:
        print(f"ОШИБКА: Не удалось запустить на порту 80. {e}")
        print("Попробуйте запустить скрипт от имени администратора или использовать другой порт (например, 5000).")