# Лабораторная работа 6 - Рефакторинг и Code Smells

## Описание проекта

Проект демонстрирует рефакторинг модуля обработки заказов (`order_processing.py`). 
Исходный код содержал множество проблем качества, которые были устранены в процессе рефакторинга.

## Выявленные проблемы качества (Code Smells)

### 1. **Long Method (Длинный метод)**
Функция `process_checkout` выполняла слишком много обязанностей:
- Парсинг запроса
- Валидацию данных
- Расчёт суммы
- Применение скидок
- Расчёт налогов
- Генерацию ID
- Формирование ответа

**Проблема:** Функция из 60+ строк сложна для понимания и тестирования.

### 2. **Magic Numbers (Магические числа)**
Код содержал множество магических чисел без объяснения их значения:
- `0.10`, `0.20`, `0.05` - размеры скидок
- `0.21` - ставка налога
- `200` - порог для купона SAVE20
- `100` - порог для купона VIP
- `50`, `10` - фиксированные скидки VIP

**Проблема:** Непонятно, что означают числа, сложно изменять.

### 3. **Complex Conditional (Сложная условная логика)**
Расчёт скидок представлял собой длинную цепочку `if-elif-else`:
```python
if coupon is None or coupon == "":
    discount = 0
elif coupon == "SAVE10":
    discount = int(subtotal * 0.10)
elif coupon == "SAVE20":
    if subtotal >= 200:
        discount = int(subtotal * 0.20)
    else:
        discount = int(subtotal * 0.05)
# ... и так далее
```

**Проблема:** Сложно читать, добавлять новые купоны, тестировать отдельно.

### 4. **Duplicated Code (Дублирование кода)**
Повторяющиеся паттерны:
- Проверки на `None` и пустые значения
- Валидация каждого поля товара в цикле
- Конвертация в `int` при расчёте скидок

### 5. **Primitive Obsession (Одержимость примитивами)**
Использование словарей вместо структурированных объектов:
```python
user_id, items, coupon, currency = parse_request(request)
```

**Проблема:** Нет типизации, легко перепутать порядок параметров.

### 6. **Poor Naming (Плохие имена)**
Неинформативные имена переменных:
- `it` вместо `item`
- `r` в тестах вместо `result` или `order`

### 7. **Shotgun Surgery (Хирургия дробью)**
Добавление нового купона требовало изменений в нескольких местах:
- Добавление условия в `if-elif`
- Добавление магических чисел
- Обновление тестов

### 8. **Feature Envy (Зависть к функциональности)**
Функция `process_checkout` сама выполняла все вычисления вместо делегирования.

### 9. **Low Cohesion (Низкая связность)**
Валидация, бизнес-логика и форматирование ответа смешаны в одной функции.

## Применённые рефакторинги

### 1. **Extract Method (Извлечение метода)**
Большая функция разбита на специализированные функции:
- `parse_request()` - парсинг запроса
- `validate_request()` - валидация
- `calculate_subtotal()` - расчёт промежуточной суммы
- `calculate_discount()` - расчёт скидки
- `calculate_tax()` - расчёт налога
- `apply_discount()` - применение скидки
- `generate_order_id()` - генерация ID
- `build_order_response()` - формирование ответа

### 2. **Extract Constant (Извлечение константы)**
Все магические числа заменены на именованные константы:
```python
DEFAULT_CURRENCY = "USD"
TAX_RATE = 0.21
COUPON_SAVE10_DISCOUNT = 0.10
COUPON_SAVE20_DISCOUNT = 0.20
COUPON_SAVE20_THRESHOLD = 200
# ... и т.д.
```

### 3. **Introduce Parameter Object (Введение объекта-параметра)**
Созданы dataclass'ы для структурирования данных:
```python
@dataclass
class OrderRequest:
    user_id: int
    items: List[Dict[str, int]]
    coupon: Optional[str]
    currency: str
```

### 4. **Replace Conditional with Polymorphism (Замена условий полиморфизмом)**
Вместо длинной цепочки `if-elif` для купонов, создали отдельные функции:
```python
def _calculate_save10_discount(subtotal: int) -> int:
    return int(subtotal * COUPON_SAVE10_DISCOUNT)

def _calculate_save20_discount(subtotal: int) -> int:
    if subtotal >= COUPON_SAVE20_THRESHOLD:
        return int(subtotal * COUPON_SAVE20_DISCOUNT)
    return int(subtotal * COUPON_SAVE20_FALLBACK_DISCOUNT)
```

### 5. **Decompose Conditional (Декомпозиция условий)**
Сложные условия вынесены в отдельные функции с понятными именами.

### 6. **Add Type Hints (Добавление аннотаций типов)**
Все функции получили аннотации типов для параметров и возвращаемых значений:
```python
def calculate_subtotal(items: List[Dict[str, int]]) -> int:
```

### 7. **Add Docstrings (Добавление документации)**
Каждая функция документирована в формате docstring:
```python
def calculate_discount(subtotal: int, coupon: Optional[str]) -> int:
    """
    Расчёт скидки по купону.
    
    Args:
        subtotal: Промежуточная сумма заказа
        coupon: Код купона
        
    Returns:
        int: Размер скидки
        
    Raises:
        ValueError: Если купон неизвестен
    """
```

### 8. **Rename Variables (Переименование переменных)**
Все переменные получили понятные имена:
- `it` → `item`
- `r` → `order` или `result`

### 9. **Single Responsibility Principle (Принцип единственной ответственности)**
Каждая функция теперь отвечает за одну конкретную задачу.

## Результаты рефакторинга

### Что стало проще читать:

1. **Главная функция как сценарий**
   ```python
   def process_checkout(request: dict) -> dict:
       order_request = parse_request(request)
       validate_request(order_request)
       subtotal = calculate_subtotal(order_request.items)
       discount = calculate_discount(subtotal, order_request.coupon)
       total_after_discount = apply_discount(subtotal, discount)
       tax = calculate_tax(total_after_discount)
       total = total_after_discount + tax
       order_id = generate_order_id(order_request.user_id, len(order_request.items))
       return build_order_response(...)
   ```
   Теперь читается как пошаговый алгоритм!

2. **Самодокументируемый код**
   - Имена функций объясняют, что они делают
   - Константы объясняют значения
   - Типы помогают понять структуру данных

3. **Изолированная логика**
   - Каждую функцию можно понять независимо
   - Меньше когнитивной нагрузки

### Что стало проще менять:

1. **Добавление нового купона**
   - Было: изменения в 3-4 местах
   - Стало: добавить константы и одну функцию `_calculate_XXX_discount`

2. **Изменение ставки налога**
   - Было: поиск магического числа `0.21` по коду
   - Стало: изменить константу `TAX_RATE`

3. **Изменение логики валидации**
   - Было: править внутри огромной функции
   - Стало: править функцию `validate_request`

4. **Тестирование**
   - Каждую функцию можно тестировать отдельно
   - Легче создавать unit-тесты для специфичной логики

### Метрики улучшения:

| Метрика | До | После |
|---------|-----|-------|
| Строк в `process_checkout` | 60+ | 20 |
| Количество функций | 2 | 13 |
| Cyclomatic Complexity | ~15 | ~3 в каждой функции |
| Магических чисел | 8 | 0 |
| Максимальная вложенность | 3 | 1-2 |

## Принципы, применённые в коде

### DRY (Don't Repeat Yourself)
- Повторяющиеся проверки вынесены в функции
- Логика расчёта скидок не дублируется

### KISS (Keep It Simple, Stupid)
- Каждая функция делает одну простую вещь
- Нет сложной вложенной логики

### SOLID
- **S (Single Responsibility)**: каждая функция имеет одну ответственность
- **O (Open/Closed)**: легко добавить новый тип купона без изменения существующего кода
- **D (Dependency Inversion)**: функции зависят от абстракций (типов), а не от конкретики

## Как запустить

### Установка зависимостей
```bash
pip install pytest
```

### Запуск тестов
```bash
# Для исходного кода
pytest test_order_processing.py -v

# После рефакторинга (измените import в тестах)
pytest test_order_processing.py -v
```

### Проверка покрытия
```bash
pip install pytest-cov
pytest test_order_processing.py --cov=order_processing_refactored --cov-report=html
```

## Выводы

Рефакторинг значительно улучшил качество кода:
- ✅ Код стал читаемым и самодокументируемым
- ✅ Уменьшилась сложность функций
- ✅ Появилась возможность переиспользования компонентов
- ✅ Упростилось тестирование
- ✅ Снизились риски при внесении изменений
- ✅ Код соответствует принципам SOLID, DRY, KISS

Все изменения были сделаны маленькими шагами с постоянным запуском тестов, 
что гарантировало сохранение поведения системы.