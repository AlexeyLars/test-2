package domain

import "errors"

// OrderLine - строка заказа (часть агрегата Order)
type OrderLine struct {
	productID string
	price     Money
	quantity  int
}

// NewOrderLine создаёт новую строку заказа
func NewOrderLine(productID string, price Money, quantity int) (OrderLine, error) {
	if productID == "" {
		return OrderLine{}, errors.New("productID cannot be empty")
	}
	if quantity <= 0 {
		return OrderLine{}, errors.New("quantity must be positive")
	}
	return OrderLine{
		productID: productID,
		price:     price,
		quantity:  quantity,
	}, nil
}

// ProductID возвращает ID продукта
func (ol OrderLine) ProductID() string {
	return ol.productID
}

// Price возвращает цену за единицу
func (ol OrderLine) Price() Money {
	return ol.price
}

// Quantity возвращает количество
func (ol OrderLine) Quantity() int {
	return ol.quantity
}

// Total рассчитывает общую стоимость строки
func (ol OrderLine) Total() Money {
	totalAmount := ol.price.Amount() * int64(ol.quantity)
	// Ошибка не возникнет, т.к. мы используем ту же валюту
	total, _ := NewMoney(totalAmount, ol.price.Currency())
	return total
}
