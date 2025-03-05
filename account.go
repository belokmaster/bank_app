package main

import (
	"errors"
	"fmt"
)

// Account представляет собой банковский счёт
type Account struct {
	balance float64
}

// NewAccount создаёт новый счёт с начальным балансом
func NewAccount(initialBalance float64) *Account {
	return &Account{balance: initialBalance}
}

// Deposit добавляет деньги на счёт
func (a *Account) Deposit(amount float64) {
	a.balance += amount
	fmt.Printf("Добавлено: %.2f, Новый баланс: %.2f\n", amount, a.balance)
}

// Withdraw снимает деньги со счёта, если хватает средств
func (a *Account) Withdraw(amount float64) error {
	if amount > a.balance {
		return errors.New("недостаточно средств")
	}

	a.balance -= amount
	fmt.Printf("Снято: %.2f, Новый баланс: %.2f\n", amount, a.balance)
	return nil
}

// GetBalance возвращает текущий баланс
func (a *Account) GetBalance() float64 {
	return a.balance
}
