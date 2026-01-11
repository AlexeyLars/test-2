package domain

import (
	"errors"
	"fmt"
)

// Money - value object для представления денежной суммы
type Money struct {
	amount   int64  // сумма в минимальных единицах (например, копейки)
	currency string
}

// NewMoney создаёт новый объект Money
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{amount: amount, currency: currency}, nil
}

// Amount возвращает сумму
func (m Money) Amount() int64 {
	return m.amount
}

// Currency возвращает валюту
func (m Money) Currency() string {
	return m.currency
}

// Add складывает две суммы
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Equals проверяет равенство двух сумм
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// IsZero проверяет, равна ли сумма нулю
func (m Money) IsZero() bool {
	return m.amount == 0
}

// String возвращает строковое представление
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.amount)/100.0, m.currency)
}
