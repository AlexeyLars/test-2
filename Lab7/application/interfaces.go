package application

import "lab7/domain"

// OrderRepository - интерфейс для работы с хранилищем заказов
type OrderRepository interface {
	// GetByID загружает заказ по идентификатору
	GetByID(orderID string) (*domain.Order, error)

	// Save сохраняет заказ
	Save(order *domain.Order) error
}

// PaymentGateway - интерфейс для проведения платежей
type PaymentGateway interface {
	// Charge выполняет списание средств
	Charge(orderID string, money domain.Money) error
}
