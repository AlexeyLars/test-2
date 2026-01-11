package infrastructure

import (
	"errors"
	"lab7/domain"
	"sync"
)

// InMemoryOrderRepository - реализация OrderRepository в памяти
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

// NewInMemoryOrderRepository создаёт новый репозиторий в памяти
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

// GetByID загружает заказ по идентификатору
func (r *InMemoryOrderRepository) GetByID(orderID string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}

	// Возвращаем копию заказа для изоляции
	return r.copyOrder(order), nil
}

// Save сохраняет заказ
func (r *InMemoryOrderRepository) Save(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Сохраняем копию заказа
	r.orders[order.ID()] = r.copyOrder(order)
	return nil
}

// copyOrder создаёт копию заказа для изоляции
func (r *InMemoryOrderRepository) copyOrder(order *domain.Order) *domain.Order {
	lines := order.Lines()
	return domain.ReconstructOrder(order.ID(), lines, order.Status())
}
