# Лаба 7 - Архитектура, слои и DDD-lite

На лекции мы разобрали слоистую архитектуру, DIP и основы DDD-lite.
В этой лабораторной работе реализована система оплаты заказа с разделением по слоям и доменной моделью.

## Структура проекта

```
Lab7/
├── domain/                    # Доменный слой
│   ├── money.go              # Value Object для денежных сумм
│   ├── order.go              # Агрегат заказа
│   ├── order_line.go         # Строка заказа (часть агрегата)
│   └── order_status.go       # Статусы заказа
├── application/               # Слой приложения
│   ├── interfaces.go         # Интерфейсы для репозитория и платёжного шлюза
│   └── pay_order_use_case.go # Use-case оплаты заказа
├── infrastructure/            # Инфраструктурный слой
│   ├── in_memory_order_repository.go  # In-memory реализация репозитория
│   └── fake_payment_gateway.go        # Фейковый платёжный шлюз
├── tests/                     # Тесты
│   └── pay_order_use_case_test.go
├── main.go                    # Пример использования
├── go.mod
└── README.md
```

## Реализация

### 1. Domain (доменный слой)

#### Money (Value Object)
- Представляет денежную сумму с валютой
- Неизменяемый объект
- Поддерживает сложение сумм в одной валюте
- Хранит сумму в минимальных единицах (копейки)

#### OrderStatus (перечисление)
- `PENDING` - заказ создан, не оплачен
- `PAID` - заказ оплачен

#### OrderLine (часть агрегата)
- Содержит информацию о товаре: ID, цена, количество
- Вычисляет общую стоимость строки

#### Order (агрегат)
- Корневая сущность агрегата
- Управляет коллекцией строк заказа
- Реализует бизнес-правила (инварианты):
  - ✅ нельзя оплатить пустой заказ
  - ✅ нельзя оплатить заказ повторно
  - ✅ после оплаты нельзя менять строки заказа
  - ✅ итоговая сумма равна сумме строк

### 2. Application (слой приложения)

#### Интерфейсы

**OrderRepository**
```go
type OrderRepository interface {
    GetByID(orderID string) (*domain.Order, error)
    Save(order *domain.Order) error
}
```

**PaymentGateway**
```go
type PaymentGateway interface {
    Charge(orderID string, money domain.Money) error
}
```

#### PayOrderUseCase
Use-case оплаты заказа, который:
1. Загружает заказ через `OrderRepository`
2. Выполняет доменную операцию оплаты (`Order.Pay()`)
3. Рассчитывает итоговую сумму
4. Вызывает платёж через `PaymentGateway`
5. Сохраняет обновлённый заказ
6. Возвращает результат оплаты

### 3. Infrastructure (инфраструктурный слой)

#### InMemoryOrderRepository
- Хранит заказы в памяти (`map[string]*Order`)
- Thread-safe реализация с использованием `sync.RWMutex`
- Создаёт копии заказов для изоляции

#### FakePaymentGateway
- Фейковая реализация для тестирования
- Записывает все платежи в память
- Может симулировать ошибки
- Предоставляет методы для проверки платежей в тестах

### 4. Tests (тесты)

Реализованы все требуемые тесты:

✅ **TestPayOrder_Success** - успешная оплата корректного заказа
✅ **TestPayOrder_EmptyOrder** - ошибка при оплате пустого заказа
✅ **TestPayOrder_AlreadyPaid** - ошибка при повторной оплате
✅ **TestOrder_CannotModifyAfterPayment** - невозможность изменения заказа после оплаты
✅ **TestOrder_TotalCalculation** - корректный расчёт итоговой суммы

Дополнительные тесты:
- **TestPayOrder_PaymentGatewayFailure** - обработка ошибки платёжного шлюза
- **TestOrder_DifferentCurrencies** - проверка ошибки при смешивании валют

## Запуск

### Запуск тестов
```bash
go test ./tests/... -v
```

### Запуск примера
```bash
go run main.go
```

## Пример использования

```go
// Создаём инфраструктуру
repo := infrastructure.NewInMemoryOrderRepository()
gateway := infrastructure.NewFakePaymentGateway()

// Создаём use-case
payOrderUseCase := application.NewPayOrderUseCase(repo, gateway)

// Создаём заказ
order := domain.NewOrder("order-123")

// Добавляем товары
price1, _ := domain.NewMoney(15000, "RUB") // 150.00 RUB
line1, _ := domain.NewOrderLine("laptop", price1, 1)
order.AddLine(line1)

price2, _ := domain.NewMoney(5000, "RUB") // 50.00 RUB
line2, _ := domain.NewOrderLine("mouse", price2, 2)
order.AddLine(line2)

// Сохраняем заказ
repo.Save(order)

// Оплачиваем заказ
result, err := payOrderUseCase.Execute("order-123")
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Message)
// Output: order order-123 paid successfully for 250.00 RUB
```

## Архитектурные принципы

### Dependency Inversion Principle (DIP)
- Слой `application` зависит от интерфейсов, а не от конкретных реализаций
- Слой `infrastructure` реализует интерфейсы из слоя `application`
- Доменный слой не зависит ни от чего

### Слоистая архитектура
```
Tests → Application → Domain
           ↓
    Infrastructure
```

- **Domain** - бизнес-логика и правила, не зависит от других слоёв
- **Application** - оркеструет use-case, зависит только от Domain
- **Infrastructure** - технические детали, зависит от Application и Domain
- **Tests** - тестирует use-case без реальной базы данных

### DDD-lite принципы
- **Aggregate** - Order управляет своими OrderLine
- **Value Object** - Money - неизменяемый объект
- **Entity** - Order имеет идентификатор
- **Invariants** - бизнес-правила защищены в доменной модели
- **Repository Pattern** - абстракция для работы с хранилищем

## Результаты тестирования

```
=== RUN   TestPayOrder_Success
--- PASS: TestPayOrder_Success (0.00s)
=== RUN   TestPayOrder_EmptyOrder
--- PASS: TestPayOrder_EmptyOrder (0.00s)
=== RUN   TestPayOrder_AlreadyPaid
--- PASS: TestPayOrder_AlreadyPaid (0.00s)
=== RUN   TestOrder_CannotModifyAfterPayment
--- PASS: TestOrder_CannotModifyAfterPayment (0.00s)
=== RUN   TestOrder_TotalCalculation
--- PASS: TestOrder_TotalCalculation (0.00s)
=== RUN   TestPayOrder_PaymentGatewayFailure
--- PASS: TestPayOrder_PaymentGatewayFailure (0.00s)
=== RUN   TestOrder_DifferentCurrencies
--- PASS: TestOrder_DifferentCurrencies (0.00s)
PASS
```

Все тесты пройдены успешно! ✅
