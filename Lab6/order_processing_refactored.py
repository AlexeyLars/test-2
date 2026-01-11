"""
Модуль обработки заказов с применением рефакторинга.
"""

from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass

# Константы
DEFAULT_CURRENCY = "USD"
TAX_RATE = 0.21

# Купоны и их параметры
COUPON_SAVE10_DISCOUNT = 0.10
COUPON_SAVE20_DISCOUNT = 0.20
COUPON_SAVE20_FALLBACK_DISCOUNT = 0.05
COUPON_SAVE20_THRESHOLD = 200
COUPON_VIP_HIGH_DISCOUNT = 50
COUPON_VIP_LOW_DISCOUNT = 10
COUPON_VIP_THRESHOLD = 100


@dataclass
class OrderRequest:
    """Структура запроса на создание заказа."""
    user_id: int
    items: List[Dict[str, int]]
    coupon: Optional[str]
    currency: str


@dataclass
class OrderItem:
    """Структура товара в заказе."""
    price: int
    qty: int


def parse_request(request: dict) -> OrderRequest:
    """
    Парсинг и нормализация входного запроса.
    
    Args:
        request: Словарь с данными запроса
        
    Returns:
        OrderRequest: Структурированный объект запроса
    """
    return OrderRequest(
        user_id=request.get("user_id"),
        items=request.get("items"),
        coupon=request.get("coupon"),
        currency=request.get("currency") or DEFAULT_CURRENCY
    )


def validate_request(order_request: OrderRequest) -> None:
    """
    Валидация данных запроса.
    
    Args:
        order_request: Объект запроса для валидации
        
    Raises:
        ValueError: Если данные невалидны
    """
    if order_request.user_id is None:
        raise ValueError("user_id is required")
    
    if order_request.items is None:
        raise ValueError("items is required")
    
    if not isinstance(order_request.items, list):
        raise ValueError("items must be a list")
    
    if len(order_request.items) == 0:
        raise ValueError("items must not be empty")
    
    _validate_items(order_request.items)


def _validate_items(items: List[Dict[str, int]]) -> None:
    """
    Валидация списка товаров.
    
    Args:
        items: Список товаров для валидации
        
    Raises:
        ValueError: Если товары невалидны
    """
    for item in items:
        if "price" not in item or "qty" not in item:
            raise ValueError("item must have price and qty")
        
        if item["price"] <= 0:
            raise ValueError("price must be positive")
        
        if item["qty"] <= 0:
            raise ValueError("qty must be positive")


def calculate_subtotal(items: List[Dict[str, int]]) -> int:
    """
    Расчёт промежуточной суммы заказа.
    
    Args:
        items: Список товаров
        
    Returns:
        int: Промежуточная сумма
    """
    return sum(item["price"] * item["qty"] for item in items)


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
    if not coupon or coupon == "":
        return 0
    
    if coupon == "SAVE10":
        return _calculate_save10_discount(subtotal)
    
    if coupon == "SAVE20":
        return _calculate_save20_discount(subtotal)
    
    if coupon == "VIP":
        return _calculate_vip_discount(subtotal)
    
    raise ValueError("unknown coupon")


def _calculate_save10_discount(subtotal: int) -> int:
    """Расчёт скидки для купона SAVE10."""
    return int(subtotal * COUPON_SAVE10_DISCOUNT)


def _calculate_save20_discount(subtotal: int) -> int:
    """Расчёт скидки для купона SAVE20."""
    if subtotal >= COUPON_SAVE20_THRESHOLD:
        return int(subtotal * COUPON_SAVE20_DISCOUNT)
    return int(subtotal * COUPON_SAVE20_FALLBACK_DISCOUNT)


def _calculate_vip_discount(subtotal: int) -> int:
    """Расчёт скидки для купона VIP."""
    if subtotal >= COUPON_VIP_THRESHOLD:
        return COUPON_VIP_HIGH_DISCOUNT
    return COUPON_VIP_LOW_DISCOUNT


def calculate_tax(amount: int) -> int:
    """
    Расчёт налога.
    
    Args:
        amount: Сумма для расчёта налога
        
    Returns:
        int: Размер налога
    """
    return int(amount * TAX_RATE)


def apply_discount(subtotal: int, discount: int) -> int:
    """
    Применение скидки к промежуточной сумме.
    
    Args:
        subtotal: Промежуточная сумма
        discount: Размер скидки
        
    Returns:
        int: Сумма после применения скидки (не меньше 0)
    """
    result = subtotal - discount
    return max(0, result)


def generate_order_id(user_id: int, items_count: int) -> str:
    """
    Генерация идентификатора заказа.
    
    Args:
        user_id: ID пользователя
        items_count: Количество товаров
        
    Returns:
        str: Идентификатор заказа
    """
    return f"{user_id}-{items_count}-X"


def build_order_response(
    order_id: str,
    user_id: int,
    currency: str,
    subtotal: int,
    discount: int,
    tax: int,
    total: int,
    items_count: int
) -> Dict:
    """
    Формирование ответа с данными заказа.
    
    Args:
        order_id: ID заказа
        user_id: ID пользователя
        currency: Валюта
        subtotal: Промежуточная сумма
        discount: Скидка
        tax: Налог
        total: Итоговая сумма
        items_count: Количество товаров
        
    Returns:
        Dict: Словарь с данными заказа
    """
    return {
        "order_id": order_id,
        "user_id": user_id,
        "currency": currency,
        "subtotal": subtotal,
        "discount": discount,
        "tax": tax,
        "total": total,
        "items_count": items_count,
    }


def process_checkout(request: dict) -> dict:
    """
    Основная функция обработки оформления заказа.
    
    Args:
        request: Словарь с данными запроса
        
    Returns:
        Dict: Результат обработки заказа
        
    Raises:
        ValueError: Если данные запроса невалидны
    """
    # 1. Парсинг запроса
    order_request = parse_request(request)
    
    # 2. Валидация
    validate_request(order_request)
    
    # 3. Расчёт промежуточной суммы
    subtotal = calculate_subtotal(order_request.items)
    
    # 4. Расчёт скидки
    discount = calculate_discount(subtotal, order_request.coupon)
    
    # 5. Применение скидки
    total_after_discount = apply_discount(subtotal, discount)
    
    # 6. Расчёт налога
    tax = calculate_tax(total_after_discount)
    
    # 7. Расчёт итоговой суммы
    total = total_after_discount + tax
    
    # 8. Генерация ID заказа
    order_id = generate_order_id(order_request.user_id, len(order_request.items))
    
    # 9. Формирование ответа
    return build_order_response(
        order_id=order_id,
        user_id=order_request.user_id,
        currency=order_request.currency,
        subtotal=subtotal,
        discount=discount,
        tax=tax,
        total=total,
        items_count=len(order_request.items)
    )
