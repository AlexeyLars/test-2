package application

import (
	"fmt"
)

// PayOrderResult - результат выполнения use-case оплаты заказа
type PayOrderResult struct {
	Success bool
	Message string
}

// PayOrderUseCase - use-case для оплаты заказа
type PayOrderUseCase struct {
	orderRepo      OrderRepository
	paymentGateway PaymentGateway
}

// NewPayOrderUseCase создаёт новый use-case
func NewPayOrderUseCase(orderRepo OrderRepository, paymentGateway PaymentGateway) *PayOrderUseCase {
	return &PayOrderUseCase{
		orderRepo:      orderRepo,
		paymentGateway: paymentGateway,
	}
}

// Execute выполняет оплату заказа
func (uc *PayOrderUseCase) Execute(orderID string) (PayOrderResult, error) {
	// 1. Загружаем заказ через OrderRepository
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return PayOrderResult{
			Success: false,
			Message: fmt.Sprintf("failed to load order: %v", err),
		}, err
	}

	// 2. Выполняем доменную операцию оплаты
	err = order.Pay()
	if err != nil {
		return PayOrderResult{
			Success: false,
			Message: fmt.Sprintf("failed to pay order: %v", err),
		}, err
	}

	// 3. Рассчитываем сумму для оплаты
	total, err := order.Total()
	if err != nil {
		return PayOrderResult{
			Success: false,
			Message: fmt.Sprintf("failed to calculate total: %v", err),
		}, err
	}

	// 4. Вызываем платёж через PaymentGateway
	err = uc.paymentGateway.Charge(orderID, total)
	if err != nil {
		return PayOrderResult{
			Success: false,
			Message: fmt.Sprintf("payment failed: %v", err),
		}, err
	}

	// 5. Сохраняем заказ
	err = uc.orderRepo.Save(order)
	if err != nil {
		return PayOrderResult{
			Success: false,
			Message: fmt.Sprintf("failed to save order: %v", err),
		}, err
	}

	// 6. Возвращаем результат оплаты
	return PayOrderResult{
		Success: true,
		Message: fmt.Sprintf("order %s paid successfully for %s", orderID, total.String()),
	}, nil
}
