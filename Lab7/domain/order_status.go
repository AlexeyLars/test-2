package domain

// OrderStatus - статус заказа
type OrderStatus string

const (
	// OrderStatusPending - заказ создан, но не оплачен
	OrderStatusPending OrderStatus = "PENDING"
	// OrderStatusPaid - заказ оплачен
	OrderStatusPaid OrderStatus = "PAID"
)

// String возвращает строковое представление статуса
func (s OrderStatus) String() string {
	return string(s)
}
