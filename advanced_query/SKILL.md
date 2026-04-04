---
name: gorm-advanced-query
description: Use when building complex queries with GORM including subqueries, locking, batch processing, group conditions, and scopes.
---

# Advanced Query

## Smart Select Fields

GORM can automatically select specific fields when scanning into a struct with fewer fields than the model.

```go
type APIUser struct {
  ID   uint
  Name string
}

// GORM automatically selects only `id`, `name` fields
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SQL: SELECT `id`, `name` FROM `users` LIMIT 10
```

## Locking

GORM supports different types of locks for transaction safety.

```go
// FOR UPDATE lock - prevents other transactions from modifying rows
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// SQL: SELECT * FROM `users` FOR UPDATE

// FOR SHARE lock - allows reads but prevents updates
db.Clauses(clause.Locking{Strength: "SHARE"}).Find(&users)

// NOWAIT option - fail immediately if lock unavailable
db.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Find(&users)

// SKIP LOCKED - skip locked rows in high concurrency
db.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Find(&users)
```

## SubQuery

GORM generates subqueries automatically when using a `*gorm.DB` object as a parameter.

```go
// Simple subquery
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
// SQL: SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

// FROM subquery
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&User{})
// SQL: SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18
```

## Group Conditions

Build complex queries with nested conditions using chained `Where` and `Or`.

```go
db.Where(
  db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
).Or(
  db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
).Find(&Pizza{})
// SQL: SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
```

## FirstOrInit / FirstOrCreate

```go
// FirstOrInit - fetch or initialize (no database write)
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)

// FirstOrCreate - fetch or create (writes to database)
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
```

## FindInBatches

Process large datasets efficiently in batches.

```go
result := db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
  for _, result := range results {
    // Process each record
  }
  tx.Save(&results)
  return nil // Return error to stop processing
})
```

## Scopes

Define reusable query conditions.

```go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
  return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

// Use scopes
db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
```

## Count

```go
var count int64
db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)

// Count with Distinct
db.Model(&User{}).Distinct("name").Count(&count)
```

## When NOT to Use

- **Simple CRUD operations** - Stick to basic `Create`, `First`, `Find`, `Update`, `Delete` for clarity
- **When performance is critical and queries are simple** - Raw SQL can be more performant as it avoids GORM's overhead
- **Extremely complex analytical queries** - For deep analysis with multiple CTEs, window functions, and temporary tables, raw SQL is often more suitable
- **When you don't understand the generated SQL** - If the query becomes too complex to reason about, simplify it or use raw SQL
- **Locking without a transaction** - `FOR UPDATE` and `FOR SHARE` locks should be used within an explicit transaction to be effective

## Pluck

Query a single column into a slice.

```go
var ages []int64
db.Model(&User{}).Pluck("age", &ages)

var names []string
db.Model(&User{}).Distinct().Pluck("name", &names)
```
