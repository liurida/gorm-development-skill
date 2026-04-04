
# Key Concepts for GORM Method Chaining

This document provides key concepts for using method chaining in GORM.

## Overview

GORM's API is designed to be chainable, allowing you to build complex queries by linking methods together in a fluent and readable way.

## Method Categories

GORM methods can be categorized into three types:

1.  **Chainable Methods**: These methods build up a query by adding clauses (e.g., `Where`, `Select`, `Order`, `Limit`, `Offset`, `Joins`, `Preload`). They return a `*gorm.DB` object, allowing you to chain more methods.

2.  **Finisher Methods**: These methods execute the query and perform the database operation (e.g., `First`, `Find`, `Create`, `Update`, `Delete`, `Scan`, `Count`). They also return a `*gorm.DB` object, from which you can get the result and any errors.

3.  **New Session Methods**: These methods create a new, isolated session (e.g., `Session`, `WithContext`, `Debug`). This is crucial for goroutine safety and preventing query conditions from leaking between different logical operations.

## Building a Query

You build a query by starting with a `*gorm.DB` instance and adding chainable methods.

```go
db.Model(&User{}).Where("age > ?", 18).Order("name asc").Limit(10)
```

This chain of methods constructs a query but does not execute it. The query is executed only when a finisher method is called.

```go
var users []User
err := db.Model(&User{}).Where("age > ?", 18).Order("name asc").Limit(10).Find(&users).Error
```

## Goroutine Safety and Reusability

A key concept in GORM is that a `*gorm.DB` object is **not safe for reuse** across different logical operations or goroutines after a chainable or finisher method has been called on it. This is because the object's internal state is modified.

### The Problem: Unsafe Reuse

```go
query := db.Where("name = ?", "john")

// First query
query.Where("age = ?", 30).Find(&users1) // SQL: ... WHERE name = 'john' AND age = 30

// Second query - THIS IS WRONG!
query.Where("age = ?", 40).Find(&users2) // SQL: ... WHERE name = 'john' AND age = 30 AND age = 40
```

The second query incorrectly includes the `age = 30` condition from the first query.

### The Solution: New Session Methods

To safely reuse a base query, you must create a new session from it before adding more conditions.

```go
baseQuery := db.Where("name = ?", "john").Session(&gorm.Session{})

// First query is now isolated
baseQuery.Where("age = ?", 30).Find(&users1) // SQL: ... WHERE name = 'john' AND age = 30

// Second query is also isolated
baseQuery.Where("age = ?", 40).Find(&users2) // SQL: ... WHERE name = 'john' AND age = 40
```

Using `Session(&gorm.Session{})`, `WithContext(ctx)`, or `Debug()` creates a safely reusable `*gorm.DB` instance.

This pattern is essential for writing correct, predictable, and concurrent database code with GORM.
