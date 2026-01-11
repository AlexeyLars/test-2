# Архитектура проекта Lab7

## Диаграмма зависимостей

```
┌─────────────────────────────────────────────────────┐
│                     Tests Layer                      │
│  - pay_order_use_case_test.go                       │
└───────────────────┬─────────────────────────────────┘
                    │ использует
                    ↓
┌─────────────────────────────────────────────────────┐
│                Application Layer                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ PayOrderUseCase                              │   │
│  │  1. GetByID(orderID)                         │   │
│  │  2. order.Pay()                              │   │
│  │  3. order.Total()                            │   │
│  │  4. Charge(orderID, money)                   │   │
│  │  5. Save(order)                              │   │
│  └──────────────────────────────────────────────┘   │
│                    ↓ зависит от                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ Interfaces (DIP)                             │   │
│  │  - OrderRepository (GetByID, Save)           │   │
│  │  - PaymentGateway (Charge)                   │   │
│  └──────────────────────────────────────────────┘   │
└───────────────────┬─────────────────────────────────┘
                    │ использует
                    ↓
┌─────────────────────────────────────────────────────┐
│                   Domain Layer                       │
│  ┌──────────────────────────────────────────────┐   │
│  │ Order (Aggregate Root)                       │   │
│  │  - id: string                                │   │
│  │  - lines: []OrderLine                        │   │
│  │  - status: OrderStatus                       │   │
│  │  + Pay() error                               │   │
│  │  + AddLine(line) error                       │   │
│  │  + Total() (Money, error)                    │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │ OrderLine (Entity)                           │   │
│  │  - productID: string                         │   │
│  │  - price: Money                              │   │
│  │  - quantity: int                             │   │
│  │  + Total() Money                             │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │ Money (Value Object)                         │   │
│  │  - amount: int64                             │   │
│  │  - currency: string                          │   │
│  │  + Add(Money) (Money, error)                 │   │
│  │  + Equals(Money) bool                        │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │ OrderStatus (Enum)                           │   │
│  │  - PENDING                                   │   │
│  │  - PAID                                      │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
                    ↑
                    │ реализует интерфейсы
                    │
┌─────────────────────────────────────────────────────┐
│              Infrastructure Layer                    │
│  ┌──────────────────────────────────────────────┐   │
│  │ InMemoryOrderRepository                      │   │
│  │  - orders: map[string]*Order                 │   │
│  │  - mu: sync.RWMutex                          │   │
│  │  + GetByID(orderID) (*Order, error)          │   │
│  │  + Save(order *Order) error                  │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │ FakePaymentGateway                           │   │
│  │  - payments: []PaymentRecord                 │   │
│  │  - shouldFail: bool                          │   │
│  │  - mu: sync.RWMutex                          │   │
│  │  + Charge(orderID, money) error              │   │
│  │  + GetPayments() []PaymentRecord             │   │
│  │  + SetShouldFail(bool, string)               │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Принципы проектирования

### 1. Dependency Inversion Principle (DIP)

**Проблема**: Верхние слои не должны зависеть от деталей реализации нижних слоёв.

**Решение**: Application слой определяет интерфейсы (`OrderRepository`, `PaymentGateway`), а Infrastructure слой их реализует.

```go
// Application определяет интерфейс
type OrderRepository interface {
    GetByID(orderID string) (*domain.Order, error)
    Save(order *domain.Order) error
}

// Infrastructure реализует интерфейс
type InMemoryOrderRepository struct { ... }
func (r *InMemoryOrderRepository) GetByID(...) { ... }
func (r *InMemoryOrderRepository) Save(...) { ... }
```

### 2. Domain-Driven Design (DDD-lite)

#### Aggregate (Агрегат)
`Order` - корень агрегата, который управляет коллекцией `OrderLine`. Все изменения идут через методы Order.

```go
order.AddLine(line)  // ✅ правильно - через агрегат
order.lines[0] = line // ❌ неправильно - прямой доступ
```

#### Value Object
`Money` - неизменяемый объект, представляющий денежную сумму.

```go
money1, _ := domain.NewMoney(10000, "RUB")
money2, _ := money1.Add(other)  // создаёт новый объект
// money1 остаётся неизменным
```

#### Entity
`Order` имеет уникальный идентификатор, который определяет его на протяжении всего жизненного цикла.

#### Invariants (Инварианты)
Бизнес-правила, которые всегда должны быть истинными:

1. **Нельзя оплатить пустой заказ**
```go
func (o *Order) Pay() error {
    if len(o.lines) == 0 {
        return errors.New("cannot pay empty order")
    }
    // ...
}
```

2. **Нельзя оплатить заказ повторно**
```go
func (o *Order) Pay() error {
    if o.status == OrderStatusPaid {
        return errors.New("order is already paid")
    }
    // ...
}
```

3. **После оплаты нельзя менять строки заказа**
```go
func (o *Order) AddLine(line OrderLine) error {
    if o.status == OrderStatusPaid {
        return errors.New("cannot modify paid order")
    }
    // ...
}
```

4. **Итоговая сумма равна сумме строк**
```go
func (o *Order) Total() (Money, error) {
    total := o.lines[0].Total()
    for i := 1; i < len(o.lines); i++ {
        total, _ = total.Add(o.lines[i].Total())
    }
    return total, nil
}
```

### 3. Слоистая архитектура

#### Слои и их ответственность

1. **Domain Layer** (ядро)
   - Содержит бизнес-логику
   - Не зависит ни от чего
   - Защищает инварианты

2. **Application Layer** (оркестрация)
   - Координирует use-case
   - Зависит только от Domain
   - Определяет интерфейсы для Infrastructure

3. **Infrastructure Layer** (детали)
   - Реализует интерфейсы из Application
   - Работает с внешними системами (БД, API)
   - Зависит от Application и Domain

4. **Tests Layer**
   - Тестирует use-case
   - Использует фейковые реализации из Infrastructure
   - Не использует реальную БД

## Поток выполнения PayOrderUseCase

```
1. Test вызывает useCase.Execute("order-123")
        ↓
2. UseCase загружает Order через OrderRepository
        ↓
3. UseCase вызывает order.Pay() (доменная логика)
        ↓
4. Order проверяет инварианты и меняет статус
        ↓
5. UseCase вызывает order.Total()
        ↓
6. Order суммирует все OrderLine
        ↓
7. UseCase вызывает paymentGateway.Charge(id, total)
        ↓
8. PaymentGateway записывает платёж
        ↓
9. UseCase сохраняет Order через OrderRepository
        ↓
10. UseCase возвращает результат в Test
```

## Преимущества архитектуры

### 1. Тестируемость
- Можно тестировать use-case без реальной БД
- Можно подменить PaymentGateway на фейковый
- Доменная логика изолирована и легко тестируется

### 2. Независимость от фреймворков
- Доменная логика не зависит от фреймворков
- Можно легко поменять БД (Postgres → MongoDB)
- Можно легко поменять платёжную систему

### 3. Ясность бизнес-логики
- Бизнес-правила находятся в Domain слое
- Легко понять, что делает система
- Легко добавлять новые правила

### 4. Расширяемость
- Легко добавить новый use-case
- Легко добавить новую реализацию репозитория
- Легко добавить новые способы оплаты

## Возможные улучшения

1. **Domain Events**
   - `OrderPaidEvent` - событие оплаты заказа
   - Можно отправлять уведомления, обновлять склад и т.д.

2. **Specifications**
   - Вынести проверки в отдельные спецификации
   - `CanPayOrderSpecification`

3. **Repository Query Objects**
   - `OrdersByStatusQuery`
   - `OrdersByDateRangeQuery`

4. **Result Pattern**
   - Вместо `(value, error)` использовать `Result<Value, Error>`
   - Более явная обработка ошибок

5. **Persistence Models**
   - Разделить доменные модели и модели для БД
   - `domain.Order` ← mapper → `persistence.OrderModel`
