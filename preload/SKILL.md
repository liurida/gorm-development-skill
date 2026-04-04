---
name: gorm-preload
description: Use when preloading associations in GORM to avoid N+1 query problems and efficiently load related data.
---

# GORM Preloading (Eager Loading)

GORM's preloading feature allows you to load associated data in advance, which is a crucial technique for avoiding the N+1 query problem. This skill provides a comprehensive guide to using preloading effectively.

## Preload

You can use the `Preload` method to eagerly load associations.

### Preload All Associations

Use `clause.Associations` to preload all direct associations.

```go
import "gorm.io/gorm/clause"

db.Preload(clause.Associations).Find(&users)
```

### Preload Specific Associations

```go
// Preload a single association
db.Preload("CreditCard").Find(&users)

// Preload multiple associations
db.Preload("Orders").Preload("Profile").Find(&users)
```

## Preload with Conditions

You can specify conditions to filter the preloaded associations.

```go
// Inline conditions
db.Preload("Orders", "state = ?", "paid").Find(&users)

// Custom preloading function for more complex logic
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
    return db.Order("orders.amount DESC")
}).Find(&users)
```

## Nested Preloading

GORM supports nested preloading using dot notation.

```go
// Preload nested associations
db.Preload("Orders.OrderItems.Product").Find(&users)

// Preload with conditions on nested associations
db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)
```

## When NOT to Use

- **Single record lookups** - Preloading adds queries even for one parent record; consider lazy loading for single-record views
- **When you don't need the associated data** - Preloading loads all association data; use explicit `Select` if you only need counts or IDs
- **Deep nesting with large datasets** - Preloading 4+ levels deep can generate many queries; consider denormalization or custom queries
- **When associations have many records** - Preloading a user's 10,000 orders loads all into memory; use pagination or separate queries
- **API responses needing specific shapes** - Use `Joins` with `Select` for custom projections instead of full model preloading
- **When you need JOIN filtering on parent** - Preload doesn't filter parents based on association data; use `Joins` instead
