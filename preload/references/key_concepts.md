# Key Concepts for GORM Preloading (Eager Loading)

This document provides detailed explanations of preloading associations in GORM.

## Overview

Preloading (eager loading) allows you to load associated records in advance, solving the N+1 query problem. Instead of querying the database for each associated record individually, GORM batches these queries.

## Basic Preload

Use the `Preload` method to eagerly load associations:

```go
// Two queries executed:
// 1. SELECT * FROM users;
// 2. SELECT * FROM orders WHERE user_id IN (1,2,3,4);
db.Preload("Orders").Find(&users)
```

### Multiple Associations

Chain multiple `Preload` calls:

```go
db.Preload("Orders").Preload("Profile").Preload("Company").Find(&users)
```

## Preload All Associations

Use `clause.Associations` to preload all direct associations:

```go
import "gorm.io/gorm/clause"

db.Preload(clause.Associations).Find(&users)
```

**Important:** `clause.Associations` does NOT preload nested associations. Combine it with explicit nested preloading when needed:

```go
db.Preload("Orders.OrderItems").Preload(clause.Associations).Find(&users)
```

## Preload with Conditions

### Inline Conditions

Filter preloaded records with conditions:

```go
// Only preload orders that are not cancelled
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
```

### Custom Preloading Function

Use a function for more complex logic:

```go
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
    return db.Order("orders.price DESC").Limit(10)
}).Find(&users)
```

## Nested Preloading

Use dot notation to preload nested associations:

```go
// Loads: Users -> Orders -> OrderItems -> Product
db.Preload("Orders.OrderItems.Product").Find(&users)
```

### Conditional Nested Preloading

When filtering parent associations, GORM only preloads nested associations for matched parents:

```go
// Only preload OrderItems for paid orders
db.Preload("Orders", "state = ?", "paid").
    Preload("Orders.OrderItems").
    Find(&users)
```

## Joins Preloading vs Regular Preloading

### Regular Preload (Separate Queries)

- Uses separate queries for each association
- Works with all association types (has_one, has_many, belongs_to, many_to_many)
- Better for large result sets to avoid Cartesian product

```go
db.Preload("Orders").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);
```

### Joins Preload (LEFT JOIN)

- Uses LEFT JOIN in a single query
- Only works with one-to-one relationships (has_one, belongs_to)
- More efficient for single-record associations

```go
db.Joins("Company").Find(&users)
// SELECT users.*, Company.* FROM users LEFT JOIN companies AS Company ON ...
```

### Joins with Conditions

```go
db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)
// Only joins companies where alive = true
```

## When to Use Preload vs Joins

| Scenario | Recommended | Reason |
|----------|-------------|--------|
| has_many / many_to_many | `Preload` | Avoids Cartesian product |
| belongs_to / has_one | `Joins` | Single query, more efficient |
| Filtering by association | `Joins` | Can use in WHERE clause |
| Large collections | `Preload` | Separate queries scale better |
| Need specific columns | `Joins` + `Select` | More control over result |

## Common Patterns

### Loading User with Full Order History

```go
func GetUserWithOrders(db *gorm.DB, userID uint) (*User, error) {
    var user User
    err := db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
        return db.Order("created_at DESC")
    }).Preload("Orders.OrderItems.Product").
        First(&user, userID).Error
    return &user, err
}
```

### Combining Joins and Preload

```go
// Efficient: Joins for single records, Preload for collections
db.Joins("Company").
    Joins("Profile").
    Preload("Orders").
    Find(&users)
```

### Preload with Main Query Conditions

```go
db.Where("active = ?", true).
    Preload("Orders", "state = ?", "completed").
    Find(&users)
```

## Performance Considerations

1. **Avoid N+1 queries:** Always use Preload/Joins for associations you'll access
2. **Limit nested preloading:** Each level adds query overhead
3. **Use conditions:** Only preload what you need
4. **Consider batch size:** For very large datasets, use pagination

## Common Mistakes

### Mistake 1: Accessing Unpreloaded Associations

```go
// BAD: Orders not loaded, will be empty slice
db.First(&user, 1)
for _, order := range user.Orders { // Always empty!
    // ...
}

// GOOD: Preload first
db.Preload("Orders").First(&user, 1)
```

### Mistake 2: Using Joins for has_many

```go
// BAD: Creates Cartesian product
db.Joins("Orders").Find(&users) // Each user duplicated per order

// GOOD: Use Preload
db.Preload("Orders").Find(&users)
```

### Mistake 3: Forgetting Nested Preloading

```go
// BAD: OrderItems loaded, but Product is empty
db.Preload("Orders.OrderItems").Find(&users)

// GOOD: Include Product
db.Preload("Orders.OrderItems.Product").Find(&users)
```

## Embedded Preloading

For embedded structs, use dot notation:

```go
type Address struct {
    CountryID int
    Country   Country
}

type Org struct {
    PostalAddress   Address `gorm:"embedded;embeddedPrefix:postal_"`
    VisitingAddress Address `gorm:"embedded;embeddedPrefix:visiting_"`
}

// Preload embedded association
db.Preload("PostalAddress.Country").Find(&orgs)
```

**Note:** Embedded preloading only works with `belongs_to` associations.
