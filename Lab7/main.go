package main

import (
	"fmt"
	"lab7/application"
	"lab7/domain"
	"lab7/infrastructure"
)

func main() {
	fmt.Println("=== Lab7: Architecture, Layers, and DDD-lite ===\n")

	// Создаём инфраструктуру
	repo := infrastructure.NewInMemoryOrderRepository()
	gateway := infrastructure.NewFakePaymentGateway()

	// Создаём use-case
	payOrderUseCase := application.NewPayOrderUseCase(repo, gateway)

	// Создаём заказ
	order := domain.NewOrder("order-123")

	// Добавляем товары в заказ
	price1, _ := domain.NewMoney(15000, "RUB") // 150.00 RUB
	line1, _ := domain.NewOrderLine("laptop", price1, 1)
	order.AddLine(line1)

	price2, _ := domain.NewMoney(5000, "RUB") // 50.00 RUB
	line2, _ := domain.NewOrderLine("mouse", price2, 2)
	order.AddLine(line2)

	price3, _ := domain.NewMoney(3000, "RUB") // 30.00 RUB
	line3, _ := domain.NewOrderLine("keyboard", price3, 1)
	order.AddLine(line3)

	// Сохраняем заказ
	repo.Save(order)

	// Выводим информацию о заказе
	fmt.Printf("Order ID: %s\n", order.ID())
	fmt.Printf("Status: %s\n", order.Status())
	fmt.Printf("Lines:\n")
	for _, line := range order.Lines() {
		fmt.Printf("  - %s: %d x %s = %s\n",
			line.ProductID(),
			line.Quantity(),
			line.Price().String(),
			line.Total().String())
	}
	total, _ := order.Total()
	fmt.Printf("Total: %s\n\n", total.String())

	// Оплачиваем заказ
	fmt.Println("Paying order...")
	result, err := payOrderUseCase.Execute("order-123")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result.Message)

	// Проверяем статус заказа после оплаты
	paidOrder, _ := repo.GetByID("order-123")
	fmt.Printf("Order status after payment: %s\n", paidOrder.Status())

	// Пытаемся добавить новую строку в оплаченный заказ
	fmt.Println("\nTrying to modify paid order...")
	newPrice, _ := domain.NewMoney(2000, "RUB")
	newLine, _ := domain.NewOrderLine("cable", newPrice, 1)
	err = paidOrder.AddLine(newLine)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	// Проверяем платежи
	fmt.Println("\nPayment records:")
	payments := gateway.GetPayments()
	for i, payment := range payments {
		fmt.Printf("  %d. Order %s: %s\n", i+1, payment.OrderID, payment.Amount.String())
	}
}
