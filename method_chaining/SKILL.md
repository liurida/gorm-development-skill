---
name: gorm-method-chaining
description: Use when chaining methods in GORM to build complex queries and understand the difference between chainable, finisher, and new session methods.
---

# GORM Method Chaining

GORM's API is designed to be chainable, allowing you to build complex queries by linking methods together in a fluent and readable way.

## Method Categories

GORM methods can be categorized into three types:

1.  **Chainable Methods**: These methods build up a query by adding clauses (e.g., `Where`, `Select`, `Order`, `Limit`). They return a `*gorm.DB` object, allowing you to chain more methods.
2.  **Finisher Methods**: These methods execute the query (e.g., `First`, `Find`, `Create`, `Update`). They also return a `*gorm.DB` object.
3.  **New Session Methods**: These methods create a new, isolated session (e.g., `Session`, `WithContext`, `Debug`).

## Building a Query

You build a query by starting with a `*gorm.DB` instance and adding chainable methods. The query is executed only when a finisher method is called.

```go
var users []User
err := db.Model(&User{}).Where("age > ?", 18).Order("name asc").Limit(10).Find(&users).Error
```

## Goroutine Safety and Reusability

A `*gorm.DB` object is **not safe for reuse** across different logical operations or goroutines after a chainable or finisher method has been called on it.

### The Problem: Unsafe Reuse

```go
query := db.Where("name = ?", "john")

// First query
query.Where("age = ?", 30).Find(&users1) // SQL: ... WHERE name = 'john' AND age = 30

// Second query - THIS IS WRONG!
query.Where("age = ?", 40).Find(&users2) // SQL: ... WHERE name = 'john' AND age = 30 AND age = 40
```

### The Solution: New Session Methods

To safely reuse a base query, you must create a new session from it before adding more conditions.

```go
baseQuery := db.Where("name = ?", "john").Session(&gorm.Session{})

// First query is now isolated
baseQuery.Where("age = ?", 30).Find(&users1) // SQL: ... WHERE name = 'john' AND age = 30

// Second query is also isolated
baseQuery.Where("age = ?", 40).Find(&users2) // SQL: ... WHERE name = 'john' AND age = 40
```

## When NOT to Use

- **Reusing a query builder without creating a new session** - As shown in the example, this is a major source of bugs. Always use a new session method if you intend to branch a query.
- **Building extremely long chains** - Very long chains can become hard to read and debug. Consider breaking them down into smaller, more manageable pieces, or using Scopes for reusability.
- **When the logic is too complex for a single chain** - If your query involves complex conditional logic, building it up in a more imperative style (e.g., with if-statements) can be clearer than a single, massive chain.
