import pytest
from order_processing_refactored import process_checkout


def test_ok_no_coupon():
    r = process_checkout({"user_id": 1, "items": [{"price": 50, "qty": 2}], "coupon": None, "currency": "USD"})
    assert r["subtotal"] == 100
    assert r["discount"] == 0
    assert r["tax"] == 21
    assert r["total"] == 121


def test_ok_save10():
    r = process_checkout({"user_id": 2, "items": [{"price": 30, "qty": 3}], "coupon": "SAVE10", "currency": "USD"})
    assert r["discount"] == 9


def test_ok_save20():
    r = process_checkout({"user_id": 3, "items": [{"price": 100, "qty": 2}], "coupon": "SAVE20", "currency": "USD"})
    assert r["discount"] == 40


def test_unknown_coupon():
    with pytest.raises(ValueError):
        process_checkout({"user_id": 1, "items": [{"price": 10, "qty": 1}], "coupon": "???", "currency": "USD"})


def test_save20_below_threshold():
    r = process_checkout({"user_id": 4, "items": [{"price": 50, "qty": 2}], "coupon": "SAVE20", "currency": "USD"})
    assert r["discount"] == 5


def test_vip_above_threshold():
    r = process_checkout({"user_id": 5, "items": [{"price": 100, "qty": 1}], "coupon": "VIP", "currency": "USD"})
    assert r["discount"] == 50


def test_vip_below_threshold():
    r = process_checkout({"user_id": 6, "items": [{"price": 50, "qty": 1}], "coupon": "VIP", "currency": "USD"})
    assert r["discount"] == 10


def test_default_currency():
    r = process_checkout({"user_id": 1, "items": [{"price": 50, "qty": 2}], "coupon": None, "currency": None})
    assert r["currency"] == "USD"


def test_missing_user_id():
    with pytest.raises(ValueError, match="user_id is required"):
        process_checkout({"items": [{"price": 50, "qty": 2}], "coupon": None, "currency": "USD"})


def test_missing_items():
    with pytest.raises(ValueError, match="items is required"):
        process_checkout({"user_id": 1, "coupon": None, "currency": "USD"})


def test_empty_items():
    with pytest.raises(ValueError, match="items must not be empty"):
        process_checkout({"user_id": 1, "items": [], "coupon": None, "currency": "USD"})


def test_invalid_item_format():
    with pytest.raises(ValueError, match="item must have price and qty"):
        process_checkout({"user_id": 1, "items": [{"price": 50}], "coupon": None, "currency": "USD"})


def test_negative_price():
    with pytest.raises(ValueError, match="price must be positive"):
        process_checkout({"user_id": 1, "items": [{"price": -10, "qty": 2}], "coupon": None, "currency": "USD"})


def test_negative_qty():
    with pytest.raises(ValueError, match="qty must be positive"):
        process_checkout({"user_id": 1, "items": [{"price": 50, "qty": -1}], "coupon": None, "currency": "USD"})
