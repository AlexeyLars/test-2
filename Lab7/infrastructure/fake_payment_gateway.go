package infrastructure

import (
	"errors"
	"lab7/domain"
	"sync"
)

// PaymentRecord - запись об оплате
type PaymentRecord struct {
	OrderID string
	Amount  domain.Money
}

// FakePaymentGateway - фейковая реализация платёжного шлюза для тестирования
type FakePaymentGateway struct {
	mu            sync.RWMutex
	payments      []PaymentRecord
	shouldFail    bool
	failureReason string
}

// NewFakePaymentGateway создаёт новый фейковый платёжный шлюз
func NewFakePaymentGateway() *FakePaymentGateway {
	return &FakePaymentGateway{
		payments:   make([]PaymentRecord, 0),
		shouldFail: false,
	}
}

// Charge выполняет списание средств
func (g *FakePaymentGateway) Charge(orderID string, money domain.Money) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.shouldFail {
		return errors.New(g.failureReason)
	}

	g.payments = append(g.payments, PaymentRecord{
		OrderID: orderID,
		Amount:  money,
	})

	return nil
}

// GetPayments возвращает список всех платежей
func (g *FakePaymentGateway) GetPayments() []PaymentRecord {
	g.mu.RLock()
	defer g.mu.RUnlock()

	paymentsCopy := make([]PaymentRecord, len(g.payments))
	copy(paymentsCopy, g.payments)
	return paymentsCopy
}

// SetShouldFail устанавливает, должен ли шлюз симулировать ошибку
func (g *FakePaymentGateway) SetShouldFail(shouldFail bool, reason string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.shouldFail = shouldFail
	g.failureReason = reason
}

// Reset сбрасывает состояние шлюза
func (g *FakePaymentGateway) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.payments = make([]PaymentRecord, 0)
	g.shouldFail = false
	g.failureReason = ""
}
