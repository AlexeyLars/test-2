package domain

import (
	"errors"
)

// Order - агрегат заказа
type Order struct {
	id     string
	lines  []OrderLine
	status OrderStatus
}

// NewOrder создаёт новый заказ
func NewOrder(id string) *Order {
	return &Order{
		id:     id,
		lines:  make([]OrderLine, 0),
		status: OrderStatusPending,
	}
}

// ReconstructOrder восстанавливает заказ из хранилища
func ReconstructOrder(id string, lines []OrderLine, status OrderStatus) *Order {
	return &Order{
		id:     id,
		lines:  lines,
		status: status,
	}
}

// ID возвращает идентификатор заказа
func (o *Order) ID() string {
	return o.id
}

// Lines возвращает копию строк заказа
func (o *Order) Lines() []OrderLine {
	linesCopy := make([]OrderLine, len(o.lines))
	copy(linesCopy, o.lines)
	return linesCopy
}

// Status возвращает статус заказа
func (o *Order) Status() OrderStatus {
	return o.status
}

// AddLine добавляет строку в заказ
func (o *Order) AddLine(line OrderLine) error {
	if o.status == OrderStatusPaid {
		return errors.New("cannot modify paid order")
	}
	o.lines = append(o.lines, line)
	return nil
}

// Total рассчитывает общую стоимость заказа
func (o *Order) Total() (Money, error) {
	if len(o.lines) == 0 {
		return Money{}, errors.New("order has no lines")
	}

	// Берём валюту из первой строки
	firstLine := o.lines[0]
	total := firstLine.Total()

	// Суммируем остальные строки
	for i := 1; i < len(o.lines); i++ {
		lineTotal := o.lines[i].Total()
		var err error
		total, err = total.Add(lineTotal)
		if err != nil {
			return Money{}, err
		}
	}

	return total, nil
}

// Pay выполняет оплату заказа
func (o *Order) Pay() error {
	// Инвариант: нельзя оплатить пустой заказ
	if len(o.lines) == 0 {
		return errors.New("cannot pay empty order")
	}

	// Инвариант: нельзя оплатить заказ повторно
	if o.status == OrderStatusPaid {
		return errors.New("order is already paid")
	}

	// Проверяем, что итоговая сумма корректна
	_, err := o.Total()
	if err != nil {
		return err
	}

	o.status = OrderStatusPaid
	return nil
}

// IsPaid проверяет, оплачен ли заказ
func (o *Order) IsPaid() bool {
	return o.status == OrderStatusPaid
}
