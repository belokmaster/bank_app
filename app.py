from flask import Flask, render_template, request, redirect, url_for
import database as db
import json

app = Flask(__name__)


@app.route('/')
def index():
    transactions = db.get_all_transactions_db()
    totals = db.calculate_totals_db()

    categories_data = db.get_categories_with_subcategories_db()
    return render_template('index.html', 
                           transactions=transactions, 
                           totals=totals,
                           categories_data_json=json.dumps(categories_data))


@app.route('/add', methods=['POST'])
def add():
    trans_type = request.form['type']
    category_id = request.form['category_id']
    subcategory_id = request.form.get('subcategory_id') 
    description = request.form['description']
    
    try:
        amount = float(request.form['amount'])
        if amount <= 0:
            return "Ошибка: сумма должна быть положительным числом. <a href='/'>Назад</a>"
    except ValueError:
        return "Ошибка: сумма должна быть числом. <a href='/'>Назад</a>"

    db.add_transaction_db(trans_type, category_id, subcategory_id, amount, description)

    return redirect(url_for('index'))

@app.route('/add_category', methods=['POST'])
def add_category():
    name = request.form['name']
    type = request.form['type']

    if name and type:
        db.add_category_db(name, type)

    return redirect(url_for('index'))

@app.route('/add_subcategory', methods=['POST'])
def add_subcategory():
    name = request.form['name']
    category_id = request.form['category_id']

    if name and category_id:
        db.add_subcategory_db(name, category_id)

    return redirect(url_for('index'))

@app.route('/delete/<int:transaction_id>', methods=['POST'])
def delete(transaction_id):
    db.delete_transaction_db(transaction_id)
    return redirect(url_for('index'))


if __name__ == '__main__':
    db.init_db()  
    try:
        app.run(host='0.0.0.0', port=5000, debug=True) 
    except OSError as e:
        print(f"ОШИБКА: {e}")
        print("Попробуйте запустить скрипт от имени администратора или использовать другой порт.")