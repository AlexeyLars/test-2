package tests

import (
	"lab7/application"
	"lab7/domain"
	"lab7/infrastructure"
	"testing"
)

// setupTestEnvironment создаёт тестовое окружение
func setupTestEnvironment() (*infrastructure.InMemoryOrderRepository, *infrastructure.FakePaymentGateway, *application.PayOrderUseCase) {
	repo := infrastructure.NewInMemoryOrderRepository()
	gateway := infrastructure.NewFakePaymentGateway()
	useCase := application.NewPayOrderUseCase(repo, gateway)
	return repo, gateway, useCase
}

// TestPayOrder_Success проверяет успешную оплату корректного заказа
func TestPayOrder_Success(t *testing.T) {
	repo, gateway, useCase := setupTestEnvironment()

	// Создаём заказ с несколькими строками
	order := domain.NewOrder("order-1")
	money1, _ := domain.NewMoney(10000, "RUB") // 100.00 RUB
	line1, _ := domain.NewOrderLine("product-1", money1, 2)
	order.AddLine(line1)

	money2, _ := domain.NewMoney(5000, "RUB") // 50.00 RUB
	line2, _ := domain.NewOrderLine("product-2", money2, 1)
	order.AddLine(line2)

	// Сохраняем заказ
	repo.Save(order)

	// Выполняем оплату
	result, err := useCase.Execute("order-1")

	// Проверяем результат
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got: %v", result.Message)
	}

	// Проверяем, что заказ оплачен
	updatedOrder, _ := repo.GetByID("order-1")
	if !updatedOrder.IsPaid() {
		t.Error("expected order to be paid")
	}

	// Проверяем, что платёж прошёл
	payments := gateway.GetPayments()
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got: %d", len(payments))
	}

	expectedTotal, _ := domain.NewMoney(25000, "RUB") // 2*100 + 1*50 = 250.00 RUB
	if !payments[0].Amount.Equals(expectedTotal) {
		t.Errorf("expected payment amount %s, got: %s", expectedTotal.String(), payments[0].Amount.String())
	}
}

// TestPayOrder_EmptyOrder проверяет ошибку при оплате пустого заказа
func TestPayOrder_EmptyOrder(t *testing.T) {
	repo, _, useCase := setupTestEnvironment()

	// Создаём пустой заказ
	order := domain.NewOrder("order-2")
	repo.Save(order)

	// Пытаемся оплатить
	result, err := useCase.Execute("order-2")

	// Ожидаем ошибку
	if err == nil {
		t.Fatal("expected error for empty order, got nil")
	}
	if result.Success {
		t.Error("expected failure for empty order")
	}
}

// TestPayOrder_AlreadyPaid проверяет ошибку при повторной оплате
func TestPayOrder_AlreadyPaid(t *testing.T) {
	repo, _, useCase := setupTestEnvironment()

	// Создаём и оплачиваем заказ
	order := domain.NewOrder("order-3")
	money, _ := domain.NewMoney(10000, "RUB")
	line, _ := domain.NewOrderLine("product-1", money, 1)
	order.AddLine(line)
	order.Pay() // Оплачиваем
	repo.Save(order)

	// Пытаемся оплатить повторно
	result, err := useCase.Execute("order-3")

	// Ожидаем ошибку
	if err == nil {
		t.Fatal("expected error for already paid order, got nil")
	}
	if result.Success {
		t.Error("expected failure for already paid order")
	}
}

// TestOrder_CannotModifyAfterPayment проверяет невозможность изменения заказа после оплаты
func TestOrder_CannotModifyAfterPayment(t *testing.T) {
	// Создаём заказ
	order := domain.NewOrder("order-4")
	money, _ := domain.NewMoney(10000, "RUB")
	line, _ := domain.NewOrderLine("product-1", money, 1)
	order.AddLine(line)

	// Оплачиваем заказ
	err := order.Pay()
	if err != nil {
		t.Fatalf("expected no error when paying, got: %v", err)
	}

	// Пытаемся добавить новую строку
	newLine, _ := domain.NewOrderLine("product-2", money, 1)
	err = order.AddLine(newLine)

	// Ожидаем ошибку
	if err == nil {
		t.Fatal("expected error when modifying paid order, got nil")
	}
}

// TestOrder_TotalCalculation проверяет корректный расчёт итоговой суммы
func TestOrder_TotalCalculation(t *testing.T) {
	order := domain.NewOrder("order-5")

	// Добавляем несколько строк
	money1, _ := domain.NewMoney(15000, "RUB") // 150.00 RUB
	line1, _ := domain.NewOrderLine("product-1", money1, 3)
	order.AddLine(line1)

	money2, _ := domain.NewMoney(7500, "RUB") // 75.00 RUB
	line2, _ := domain.NewOrderLine("product-2", money2, 2)
	order.AddLine(line2)

	money3, _ := domain.NewMoney(20000, "RUB") // 200.00 RUB
	line3, _ := domain.NewOrderLine("product-3", money3, 1)
	order.AddLine(line3)

	// Рассчитываем итоговую сумму
	total, err := order.Total()
	if err != nil {
		t.Fatalf("expected no error when calculating total, got: %v", err)
	}

	// Ожидаемая сумма: 3*150 + 2*75 + 1*200 = 450 + 150 + 200 = 800.00 RUB
	expectedTotal, _ := domain.NewMoney(80000, "RUB")
	if !total.Equals(expectedTotal) {
		t.Errorf("expected total %s, got: %s", expectedTotal.String(), total.String())
	}

	// Проверяем, что сумма равна сумме строк
	var calculatedTotal domain.Money
	for i, line := range order.Lines() {
		if i == 0 {
			calculatedTotal = line.Total()
		} else {
			calculatedTotal, _ = calculatedTotal.Add(line.Total())
		}
	}

	if !total.Equals(calculatedTotal) {
		t.Errorf("total does not match sum of lines: %s vs %s", total.String(), calculatedTotal.String())
	}
}

// TestPayOrder_PaymentGatewayFailure проверяет обработку ошибки платёжного шлюза
func TestPayOrder_PaymentGatewayFailure(t *testing.T) {
	repo, gateway, useCase := setupTestEnvironment()

	// Создаём заказ
	order := domain.NewOrder("order-6")
	money, _ := domain.NewMoney(10000, "RUB")
	line, _ := domain.NewOrderLine("product-1", money, 1)
	order.AddLine(line)
	repo.Save(order)

	// Настраиваем шлюз на ошибку
	gateway.SetShouldFail(true, "insufficient funds")

	// Пытаемся оплатить
	result, err := useCase.Execute("order-6")

	// Ожидаем ошибку
	if err == nil {
		t.Fatal("expected error from payment gateway, got nil")
	}
	if result.Success {
		t.Error("expected failure when payment gateway fails")
	}
}

// TestOrder_DifferentCurrencies проверяет обработку разных валют
func TestOrder_DifferentCurrencies(t *testing.T) {
	order := domain.NewOrder("order-7")

	// Добавляем строку в RUB
	money1, _ := domain.NewMoney(10000, "RUB")
	line1, _ := domain.NewOrderLine("product-1", money1, 1)
	order.AddLine(line1)

	// Добавляем строку в USD
	money2, _ := domain.NewMoney(5000, "USD")
	line2, _ := domain.NewOrderLine("product-2", money2, 1)
	order.AddLine(line2)

	// Пытаемся рассчитать сумму
	_, err := order.Total()

	// Ожидаем ошибку о несовпадении валют
	if err == nil {
		t.Fatal("expected error for different currencies, got nil")
	}
}
