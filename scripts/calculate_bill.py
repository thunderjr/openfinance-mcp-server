#!/usr/bin/env python3
"""
Calculate the actual credit card bill amount by subtracting future installments.

Usage:
    python calculate_bill.py <balance> '<transactions_json>'

Example:
    python calculate_bill.py 1500.00 '[{"amount": 100, "total_installments": 10, "current_installment": 3}]'

The balance from the credit card includes all future installments.
This script calculates what's actually due this month by subtracting
the remaining installments that aren't due yet.

Formula:
    actual_bill = balance - SUM(amount × future_installments)
    where future_installments = total_installments - current_installment
"""

import json
import sys


def calculate_actual_bill(balance: float, transactions: list[dict]) -> float:
    """
    Calculate the actual bill amount by subtracting future installments.

    Args:
        balance: The credit card balance (includes future installments)
        transactions: List of transactions from the current bill, each with:
            - amount: The installment amount
            - total_installments: Total number of installments (optional)
            - current_installment: Current installment number (optional)

    Returns:
        The actual amount due this month
    """
    future_installments_total = 0.0

    for tx in transactions:
        total = tx.get("total_installments") or tx.get("totalInstallments")
        current = tx.get("current_installment") or tx.get("currentInstallment")
        amount = tx.get("amount", 0)

        if total is not None and current is not None:
            future = total - current
            future_installments_total += amount * future

    return balance - future_installments_total


def main():
    if len(sys.argv) != 3:
        print("Usage: python calculate_bill.py <balance> '<transactions_json>'")
        sys.exit(1)

    balance = float(sys.argv[1])
    transactions = json.loads(sys.argv[2])

    actual_bill = calculate_actual_bill(balance, transactions)

    print(json.dumps({
        "balance": balance,
        "future_installments_deducted": balance - actual_bill,
        "actual_bill": round(actual_bill, 2)
    }, indent=2))


if __name__ == "__main__":
    main()
