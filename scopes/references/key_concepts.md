
# Key Concepts for GORM Scopes

This document provides key concepts for using query scopes in GORM.

## Overview

Scopes in GORM allow you to encapsulate and reuse common query logic. A scope is a function that takes a `*gorm.DB` instance and returns a modified `*gorm.DB` instance with additional query conditions.

## Defining a Simple Scope

A basic scope is a function that adds a `Where` clause to a query.

```go
// AmountGreaterThan1000 is a scope to find orders with an amount greater than 1000.
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
    return db.Where("amount > ?", 1000)
}
```

## Using Scopes

You can apply one or more scopes to a query using the `Scopes` method.

```go
var orders []Order
// Find all orders with an amount greater than 1000
db.Scopes(AmountGreaterThan1000).Find(&orders)
```

## Scopes with Arguments

Scopes can accept arguments, allowing for more dynamic and reusable query logic.

```go
// OrderStatus is a scope to find orders with a specific status.
func OrderStatus(status string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", status)
    }
}

// Find all "paid" orders
db.Scopes(OrderStatus("paid")).Find(&orders)
```

## Chaining Scopes

You can chain multiple scopes together to build complex queries from reusable components.

```go
// PaidWithCreditCard is a scope to find orders paid by credit card.
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
    return db.Where("payment_method = ?", "credit_card")
}

// Find all paid orders with an amount > 1000 that were paid by credit card
db.Scopes(AmountGreaterThan1000, OrderStatus("paid"), PaidWithCreditCard).Find(&orders)
```

## Pagination Scope

A common use case for scopes is to handle pagination.

```go
// Paginate is a scope for paginating query results.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if page <= 0 {
            page = 1
        }
        if pageSize <= 0 {
            pageSize = 10
        }
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

// Get the second page of users
var users []User
db.Scopes(Paginate(2, 20)).Find(&users)
```

By using scopes, you can keep your query logic organized, reusable, and easy to maintain.
